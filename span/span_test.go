package span_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"

	"github.com/Strangebrewer/go-tracer/db_connection"
	"github.com/Strangebrewer/go-tracer/middleware"
	"github.com/Strangebrewer/go-tracer/server"
	"github.com/Strangebrewer/go-tracer/span"
)

const testServiceKey = "test-service-key"

var (
	testServer     *httptest.Server
	testPrivateKey *rsa.PrivateKey
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		log.Fatalf("failed to start mongo container: %v", err)
	}
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}()

	mongoURI, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	client, col, err := db_connection.Connect(ctx, mongoURI, "tracer")
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer client.Disconnect(context.Background())

	testPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate RSA key: %v", err)
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&testPrivateKey.PublicKey)
	if err != nil {
		log.Fatalf("failed to marshal public key: %v", err)
	}
	pubKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyBytes})

	authMiddleware, err := middleware.RequireAuth(string(pubKeyPEM))
	if err != nil {
		log.Fatalf("failed to create auth middleware: %v", err)
	}

	serviceKeyMiddleware := middleware.RequireServiceKey(testServiceKey)
	store := span.NewStore(col)
	srv := server.New(":0", []string{"*"}, store, authMiddleware, serviceKeyMiddleware)
	testServer = httptest.NewServer(srv.HTTPServer.Handler)
	defer testServer.Close()

	os.Exit(m.Run())
}

func makeJWT(t *testing.T, userID string) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub": userID,
		"typ": "access",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(testPrivateKey)
	require.NoError(t, err)
	return signed
}

func postSpan(t *testing.T, body any, serviceKey string) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/spans", bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if serviceKey != "" {
		req.Header.Set("X-Service-Key", serviceKey)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func TestCreateSpan_success(t *testing.T) {
	body := span.CreateSpanInput{
		TraceID:   "trace-001",
		SpanID:    "span-001",
		Service:   "go-job-search",
		Operation: "POST /applications",
		Status:    "success",
		StartTime: time.Now().Add(-50 * time.Millisecond),
		EndTime:   time.Now(),
	}

	resp := postSpan(t, body, testServiceKey)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreateSpan_unauthorized(t *testing.T) {
	body := span.CreateSpanInput{
		TraceID:   "trace-002",
		SpanID:    "span-002",
		Service:   "go-budget",
		Operation: "GET /transactions",
		Status:    "success",
		StartTime: time.Now().Add(-10 * time.Millisecond),
		EndTime:   time.Now(),
	}

	resp := postSpan(t, body, "wrong-key")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestCreateSpan_noKey(t *testing.T) {
	body := span.CreateSpanInput{
		TraceID:   "trace-003",
		SpanID:    "span-003",
		Service:   "go-auth",
		Operation: "POST /users/login",
		Status:    "success",
		StartTime: time.Now().Add(-10 * time.Millisecond),
		EndTime:   time.Now(),
	}

	resp := postSpan(t, body, "")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetTrace_success(t *testing.T) {
	traceID := "trace-get-001"
	now := time.Now().UTC().Truncate(time.Millisecond)

	// Insert two spans in reverse order to verify sort by startTime
	spans := []span.CreateSpanInput{
		{
			TraceID:   traceID,
			SpanID:    "span-b",
			Service:   "go-job-search",
			Operation: "POST /applications",
			Status:    "success",
			StartTime: now.Add(20 * time.Millisecond),
			EndTime:   now.Add(50 * time.Millisecond),
		},
		{
			TraceID:   traceID,
			SpanID:    "span-a",
			Service:   "pe-mfe-shell",
			Operation: "frontend",
			Status:    "success",
			StartTime: now,
			EndTime:   now.Add(10 * time.Millisecond),
		},
	}

	for _, s := range spans {
		resp := postSpan(t, s, testServiceKey)
		resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/traces/"+traceID, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+makeJWT(t, "user-123"))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []span.Span
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	require.Len(t, result, 2)
	assert.Equal(t, "span-a", result[0].SpanID)
	assert.Equal(t, "span-b", result[1].SpanID)
}

func TestGetTrace_unauthorized(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/traces/trace-unauth", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetTrace_empty(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/traces/trace-does-not-exist", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+makeJWT(t, "user-123"))

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result []span.Span
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.Empty(t, result)
}
