package queries

import _ "embed"

var (
	//go:embed create_link.sql
	CreateLinkQuery string
	//go:embed get_link_by_slug.sql
	GetLinkBySlugQuery string
	//go:embed increment_click_count_by_id.sql
	IncrementClickCountByIDQuery string
)
