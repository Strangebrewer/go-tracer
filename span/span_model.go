package span

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Span struct {
	ID           bson.ObjectID `bson:"_id,omitempty"          json:"-"`
	TraceID      string             `bson:"traceId"                json:"traceId"`
	SpanID       string             `bson:"spanId"                 json:"spanId"`
	ParentSpanID string             `bson:"parentSpanId,omitempty" json:"parentSpanId,omitempty"`
	Service      string             `bson:"service"                json:"service"`
	Operation    string             `bson:"operation"              json:"operation"`
	Status       string             `bson:"status"                 json:"status"`
	Error        *string            `bson:"error,omitempty"        json:"error"`
	StartTime    time.Time          `bson:"startTime"              json:"startTime"`
	EndTime      time.Time          `bson:"endTime"                json:"endTime"`
	Metadata     map[string]any     `bson:"metadata,omitempty"     json:"metadata"`
}

type CreateSpanInput struct {
	TraceID      string         `json:"traceId"`
	SpanID       string         `json:"spanId"`
	ParentSpanID string         `json:"parentSpanId,omitempty"`
	Service      string         `json:"service"`
	Operation    string         `json:"operation"`
	Status       string         `json:"status"`
	Error        *string        `json:"error,omitempty"`
	StartTime    time.Time      `json:"startTime"`
	EndTime      time.Time      `json:"endTime"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}
