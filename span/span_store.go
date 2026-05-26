package span

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store struct {
	col *mongo.Collection
}

func NewStore(col *mongo.Collection) *Store {
	return &Store{col: col}
}

func (s *Store) Create(ctx context.Context, input CreateSpanInput) error {
	doc := Span{
		TraceID:      input.TraceID,
		SpanID:       input.SpanID,
		ParentSpanID: input.ParentSpanID,
		Service:      input.Service,
		Operation:    input.Operation,
		Status:       input.Status,
		Error:        input.Error,
		StartTime:    input.StartTime,
		EndTime:      input.EndTime,
		Metadata:     input.Metadata,
	}

	_, err := s.col.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("span: failed to insert: %w", err)
	}
	return nil
}

func (s *Store) GetByTraceID(ctx context.Context, traceID string) ([]Span, error) {
	filter := bson.M{"traceId": traceID}
	opts := options.Find().SetSort(bson.D{{Key: "startTime", Value: 1}})

	cursor, err := s.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("span: failed to find: %w", err)
	}
	defer cursor.Close(ctx)

	spans := make([]Span, 0)
	if err := cursor.All(ctx, &spans); err != nil {
		return nil, fmt.Errorf("span: failed to decode: %w", err)
	}
	return spans, nil
}
