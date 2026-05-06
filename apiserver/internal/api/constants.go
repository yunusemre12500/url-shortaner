package api

import "regexp"

var SLUG_REGEX = regexp.MustCompile(`^[a-z0-9][a-z0-9-_.]{2,15}$`)
