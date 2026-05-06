package v1

import (
	"time"

	"go.yunus-emre.dev/url-shortaner/pkg/types"
)

type CreateLinkRequestBodyParams struct {
	OriginalURL string     `json:"originalUrl"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	Slug        string     `json:"slug"`
}

type CreateLinkResponseBodyParams struct {
	ClickCount  int        `json:"clickCount"`
	CreatedAt   time.Time  `json:"createdAt"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	ID          types.ID   `json:"id"`
	OriginalURL string     `json:"originalUrl"`
	Slug        string     `json:"slug"`
}
