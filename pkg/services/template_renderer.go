package services

import (
	"bytes"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/config"
	"github.com/mikestefanello/pagoda/pkg/funcs"
	"github.com/mikestefanello/pagoda/pkg/helpers"
	"github.com/mikestefanello/pagoda/pkg/page"
	"github.com/mikestefanello/pagoda/templates/layouts"
)

// TemplateRenderer provides a flexible and easy to use method of rendering Templ templates while
// also providing caching and/or hot-reloading depending on your current environment
type TemplateRenderer struct {
	// For page caching
	*pageCache

	// config stores application configuration
	config *config.Config

	// funcs stores functions to be used in templates
	funcs *funcs.Funcs
}

// NewTemplateRenderer creates a new TemplRenderer
func NewTemplateRenderer(cfg *config.Config, cache *CacheClient, web *echo.Echo) *TemplateRenderer {
	return &TemplateRenderer{
		pageCache: NewPageCache(cfg, cache),
		config:    cfg,
		funcs:     funcs.NewFuncs(web),
	}
}

// RenderPage renders a Page as an HTTP response
func (t *TemplateRenderer) RenderPage(ctx echo.Context, page *page.Page) error {
	var err error

	// Check if this is an HTMX non-boosted request which indicates that only partial
	// content should be rendered
	if page.HTMX.Request.Enabled && !page.HTMX.Request.Boosted {
		// Switch the layout which will only render the page content
		page.TemplLayout = layouts.HTMX
	}

	// If not set defaults to layouts.Main
	if page.TemplLayout == nil {
		page.TemplLayout = layouts.Main
	}

	// Set the status code
	ctx.Response().Status = page.StatusCode

	// Set any headers
	for k, v := range page.Headers {
		ctx.Response().Header().Set(k, v)
	}

	// Apply the HTMX response, if one
	if page.HTMX.Response != nil {
		page.HTMX.Response.Apply(ctx)
	}

	// Add Page and Funcs to templ context so they dont have to be passed around everywhere
	templCtx := context.WithValue(ctx.Request().Context(), helpers.TemplCtxKeyPage, page)
	templCtx = context.WithValue(templCtx, helpers.TemplCtxKeyFuncs, t.funcs)

	//Render template
	buf := &bytes.Buffer{}
	err = page.TemplLayout(page.TemplComponent).Render(templCtx, buf)
	if err != nil {
		return err
	}

	// Cache this page, if caching was enabled
	t.cachePage(ctx, page, buf)

	return ctx.HTMLBlob(ctx.Response().Status, buf.Bytes())
}
