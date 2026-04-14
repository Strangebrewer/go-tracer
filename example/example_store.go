package example

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) GetAll(ctx context.Context, userID string) ([]*Example, error) {
	// q := db.New(s.pool)
	// rows, err := q.ListExamples(ctx, userID)
	// if err != nil {
	// 	return nil, fmt.Errorf("example: GetAll: %w", err)
	// }
	// results := make([]*Example, len(rows))
	// for i, row := range rows {
	// 	results[i] = &Example{ID: row.ID, UserID: row.UserID, Name: row.Name, CreatedAt: row.CreatedAt}
	// }
	// return results, nil
	return nil, nil
}

func (s *Store) GetOne(ctx context.Context, id, userID string) (*Example, error) {
	// q := db.New(s.pool)
	// row, err := q.GetExample(ctx, db.GetExampleParams{ID: id, UserID: userID})
	// if err != nil {
	// 	return nil, fmt.Errorf("example: GetOne: %w", err)
	// }
	// return &Example{ID: row.ID, UserID: row.UserID, Name: row.Name, CreatedAt: row.CreatedAt}, nil
	return nil, nil
}

func (s *Store) Create(ctx context.Context, userID string, req CreateExampleRequest) (*Example, error) {
	// q := db.New(s.pool)
	// row, err := q.CreateExample(ctx, db.CreateExampleParams{UserID: userID, Name: req.Name})
	// if err != nil {
	// 	return nil, fmt.Errorf("example: Create: %w", err)
	// }
	// return &Example{ID: row.ID, UserID: row.UserID, Name: row.Name, CreatedAt: row.CreatedAt}, nil
	return nil, nil
}

func (s *Store) Update(ctx context.Context, id, userID string, req UpdateExampleRequest) (*Example, error) {
	// q := db.New(s.pool)
	// row, err := q.UpdateExample(ctx, db.UpdateExampleParams{ID: id, UserID: userID, Name: req.Name})
	// if err != nil {
	// 	return nil, fmt.Errorf("example: Update: %w", err)
	// }
	// return &Example{ID: row.ID, UserID: row.UserID, Name: row.Name, CreatedAt: row.CreatedAt}, nil
	return nil, nil
}

func (s *Store) Delete(ctx context.Context, id, userID string) error {
	// q := db.New(s.pool)
	// if err := q.DeleteExample(ctx, db.DeleteExampleParams{ID: id, UserID: userID}); err != nil {
	// 	return fmt.Errorf("example: Delete: %w", err)
	// }
	// return nil
	return nil
}
