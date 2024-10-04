package page

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/pkg/context"
	"github.com/mikestefanello/pagoda/pkg/htmx"
	"github.com/mikestefanello/pagoda/pkg/msg"

	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

// Page consists of all data that will be used to render a page response for a given route.
// While it's not required for a handler to render a Page on a route, this is the common data
// object that will be passed to the templates, making it easy for all handlers to share
// functionality both on the back and frontend. The Page can be expanded to include anything else
// your app wants to support.
// Methods on this page also then become available in the templates, which can be more useful than
// the funcmap if your methods require data stored in the page, such as the context.
type Page struct {
	// AppName stores the name of the application.
	// If omitted, the configuration value will be used.
	AppName string

	// Title stores the title of the page
	Title string

	// Context stores the request context
	Context echo.Context

	// Path stores the path of the current request
	Path string

	// URL stores the URL of the current request
	URL string

	// IsHome stores whether the requested page is the home page or not
	IsHome bool

	// IsAuth stores whether the user is authenticated
	IsAuth bool

	// AuthUser stores the authenticated user
	AuthUser *ent.User

	// StatusCode stores the HTTP status code that will be returned
	StatusCode int

	// Metatags stores metatag values
	Metatags struct {
		// Description stores the description metatag value
		Description string

		// Keywords stores the keywords metatag values
		Keywords []string
	}

	// Pager stores a pager which can be used to page lists of results
	Pager Pager

	// CSRF stores the CSRF token for the given request.
	// This will only be populated if the CSRF middleware is in effect for the given request.
	// If this is populated, all forms must include this value otherwise the requests will be rejected.
	CSRF string

	// Headers stores a list of HTTP headers and values to be set on the response
	Headers map[string]string

	// RequestID stores the ID of the given request.
	// This will only be populated if the request ID middleware is in effect for the given request.
	RequestID string

	// HTMX provides the ability to interact with the HTMX library
	HTMX struct {
		// Request contains the information provided by HTMX about the current request
		Request htmx.Request

		// Response contains values to pass back to HTMX
		Response *htmx.Response
	}

	// Cache stores values for caching the response of this page
	Cache struct {
		// Enabled dictates if the response of this page should be cached.
		// Cached responses are served via middleware.
		Enabled bool

		// Expiration stores the amount of time that the cache entry should live for before expiring.
		// If omitted, the configuration value will be used.
		Expiration time.Duration

		// Tags stores a list of tags to apply to the cache entry.
		// These are useful when invalidating cache for dynamic events such as entity operations.
		Tags []string
	}

	// TemplLayout stores the templ function to which TemplComponent will be passed when the page is rendered.
	// If not set defaults to layouts.Main
	TemplLayout func(templ.Component) templ.Component

	// TemplComponent stores the templ component which will be used when the page is rendered.
	// TemplRenderer will raise an error if not set
	TemplComponent templ.Component
}

// New creates and initiatizes a new Page for a given request context
func New(ctx echo.Context) *Page {
	p := &Page{
		Context:    ctx,
		Path:       ctx.Request().URL.Path,
		URL:        ctx.Request().URL.String(),
		StatusCode: http.StatusOK,
		Pager:      NewPager(ctx, DefaultItemsPerPage),
		Headers:    make(map[string]string),
		RequestID:  ctx.Response().Header().Get(echo.HeaderXRequestID),
	}

	p.IsHome = p.Path == "/"

	if csrf := ctx.Get(echomw.DefaultCSRFConfig.ContextKey); csrf != nil {
		p.CSRF = csrf.(string)
	}

	if u := ctx.Get(context.AuthenticatedUserKey); u != nil {
		p.IsAuth = true
		p.AuthUser = u.(*ent.User)
	}

	p.HTMX.Request = htmx.GetRequest(ctx)

	return p
}

// GetMessages gets all flash messages for a given type.
// This allows for easy access to flash messages from the templates.
func (p *Page) GetMessages(typ msg.Type) []string {
	return msg.Get(p.Context, typ)
}
