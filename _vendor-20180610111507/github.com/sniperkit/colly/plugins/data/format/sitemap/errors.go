package sitemap

import (
	"errors"
)

var (
	errInvalidContent           = errors.New("Could not iterate through the content submitted.")
	errInvalidSitemap           = errors.New("Invalid sitemap, could not create sitemap object.")
	errInvalidSitemapWithConfig = errors.New("Could not create sitemap object with provided config.")
)
