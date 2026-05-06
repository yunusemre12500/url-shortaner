package postgres

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.yunus-emre.dev/url-shortaner/model"
	"go.yunus-emre.dev/url-shortaner/pkg/types"
	"go.yunus-emre.dev/url-shortaner/storage"
	"go.yunus-emre.dev/url-shortaner/storage/postgres/internal/queries"
)

// ensures struct-interface compatibility at compile time
var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	config         *Config
	conn           *pgx.Conn
	connectOnce    sync.Once
	disconnectOnce sync.Once
}

func New(config *Config) (*Storage, error) {
	return &Storage{
		config: config,
	}, nil
}

// Connect implements [storage.Storage].
func (s *Storage) Connect(ctx context.Context) error {
	var err error

	s.connectOnce.Do(func() {
		s.conn, err = pgx.Connect(ctx, s.config.DSN)

		if err != nil {
			return
		}

		if err != nil {
			return
		}

		if err = s.conn.Ping(ctx); err != nil {
			return
		}
	})

	return err
}

// CreateLink implements [storage.Storage].
func (s *Storage) CreateLink(ctx context.Context, link *model.Link) error {
	if _, err := s.conn.Exec(ctx, queries.CreateLinkQuery, link.ClickCount, link.CreatedAt, link.ExpiresAt, link.ID, link.OriginalURL, link.Slug); err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == "23505" {
				return storage.ErrConflict
			}
		}
		return err
	}

	return nil
}

// Disconnect implements [storage.Storage].
func (s *Storage) Disconnect(ctx context.Context) error {
	var err error

	s.disconnectOnce.Do(func() {
		if s.conn == nil {
			return
		}

		err = s.conn.Close(ctx)
	})

	return err
}

// GetLinkBySlug implements [storage.Storage].
func (s *Storage) GetLinkBySlug(ctx context.Context, slug string) (*model.Link, error) {
	row := s.conn.QueryRow(ctx, queries.GetLinkBySlugQuery, slug)

	var link model.Link

	if err := row.Scan(&link.ClickCount, &link.CreatedAt, &link.ExpiresAt, &link.ID, &link.OriginalURL, &link.Slug); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &link, nil
}

// IncrementClickCountByLinkID implements [storage.Storage].
func (s *Storage) IncrementClickCountByLinkID(ctx context.Context, id types.ID) error {
	res, err := s.conn.Exec(ctx, queries.IncrementClickCountByIDQuery, id)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}
