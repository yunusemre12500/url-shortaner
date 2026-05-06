package v1

import (
	"errors"
	"net/url"
	"time"

	"go.yunus-emre.dev/url-shortaner/apiserver/internal/api"
)

func ValidateCreateLinkRequestBodyParams(body *CreateLinkRequestBodyParams) error {
	var errs error

	if body.ExpiresAt != nil && body.ExpiresAt.Before(time.Now()) {
		errs = errors.Join(errs, errors.New("expires as must be a future date"))
	}

	u, err := url.ParseRequestURI(body.OriginalURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		errs = errors.Join(errs, errors.New("originalUrl must be a valid absolute URL (including http/https)"))
	}

	if body.Slug != "" && !api.SLUG_REGEX.MatchString(body.Slug) {
		errs = errors.Join(errs, errors.New("slug must be between 4-17 characters and contain only lowercase alphanumeric, dots, underscores or dashes"))
	}

	return errs
}
