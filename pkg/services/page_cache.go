package services

import (
	"bytes"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/config"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/log"
	"github.com/mikestefanello/pagoda/pkg/page"
)

// cachedPageGroup stores the cache group for cached pages
const cachedPageGroup = "page"

// CachedPage is what is used to store a rendered Page in the cache
type CachedPage struct {
	// URL stores the URL of the requested page
	URL string

	// HTML stores the complete HTML of the rendered Page
	HTML []byte

	// StatusCode stores the HTTP status code
	StatusCode int

	// Headers stores the HTTP headers
	Headers map[string]string
}

type pageCache struct {
	// cache stores the cache client
	cache *CacheClient

	// config stores application configuration
	config *config.Config
}

func NewPageCache(cfg *config.Config, cache *CacheClient) *pageCache {
	return &pageCache{
		config: cfg,
		cache:  cache,
	}
}

// cachePage caches the HTML for a given Page if the Page has caching enabled
func (pc *pageCache) cachePage(ctx echo.Context, page *page.Page, html *bytes.Buffer) {
	if !page.Cache.Enabled || page.IsAuth {
		return
	}

	// If no expiration time was provided, default to the configuration value
	if page.Cache.Expiration == 0 {
		page.Cache.Expiration = pc.config.Cache.Expiration.Page
	}

	// Extract the headers
	headers := make(map[string]string)
	for k, v := range ctx.Response().Header() {
		headers[k] = v[0]
	}

	// The request URL is used as the cache key so the middleware can serve the
	// cached page on matching requests
	key := ctx.Request().URL.String()
	cp := &CachedPage{
		URL:        key,
		HTML:       html.Bytes(),
		Headers:    headers,
		StatusCode: ctx.Response().Status,
	}

	err := pc.cache.
		Set().
		Group(cachedPageGroup).
		Key(key).
		Tags(page.Cache.Tags...).
		Expiration(page.Cache.Expiration).
		Data(cp).
		Save(ctx.Request().Context())

	switch {
	case err == nil:
		log.Ctx(ctx).Debug("cached page")
	case !context.IsCanceledError(err):
		log.Ctx(ctx).Error("failed to cache page",
			"error", err,
		)
	}
}

// GetCachedPage attempts to fetch a cached page for a given URL
func (pc *pageCache) GetCachedPage(ctx echo.Context, url string) (*CachedPage, error) {
	p, err := pc.cache.
		Get().
		Group(cachedPageGroup).
		Key(url).
		Fetch(ctx.Request().Context())

	if err != nil {
		return nil, err
	}

	return p.(*CachedPage), nil
}
