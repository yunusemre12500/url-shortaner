package ulid

import (
	"github.com/oklog/ulid/v2"
	"go.yunus-emre.dev/url-shortaner/pkg/types"
)

func New() types.ID {
	return types.ID(ulid.Make().String())
}
