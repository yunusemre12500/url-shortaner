package storage

import (
	"context"
	"errors"

	"go.yunus-emre.dev/url-shortaner/model"
	"go.yunus-emre.dev/url-shortaner/pkg/types"
)

var (
	ErrConflict = errors.New("slug is conflicting with existing records")
	ErrNotFound = errors.New("not found")
)

type Storage interface {
	Connect(ctx context.Context) error
	CreateLink(ctx context.Context, link *model.Link) error
	Disconnect(ctx context.Context) error
	GetLinkBySlug(ctx context.Context, slug string) (*model.Link, error)
	IncrementClickCountByLinkID(ctx context.Context, id types.ID) error
}
