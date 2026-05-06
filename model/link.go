package model

import (
	"time"

	"go.yunus-emre.dev/url-shortaner/pkg/types"
	"go.yunus-emre.dev/url-shortaner/pkg/util/ulid"
)

type Link struct {
	ClickCount  int
	CreatedAt   time.Time
	ExpiresAt   *time.Time
	ID          types.ID
	OriginalURL string
	Slug        string
}

type CreateLinkParams struct {
	ExpiresAt   *time.Time
	OriginalURL string
	Slug        string
}

func CreateLink(params CreateLinkParams) *Link {
	return &Link{
		CreatedAt:   time.Now().UTC(),
		ExpiresAt:   params.ExpiresAt,
		ID:          ulid.New(),
		OriginalURL: params.OriginalURL,
		Slug:        params.Slug,
	}
}

func (link *Link) Expired() bool {
	return link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt)
}
