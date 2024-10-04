package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/pkg/helpers"
	"github.com/mikestefanello/pagoda/pkg/page"
	"github.com/mikestefanello/pagoda/pkg/services"
	"github.com/mikestefanello/pagoda/templates/pages"
)

const (
	routeNameAbout = "about"
	routeNameHome  = "home"
)

type (
	Pages struct {
		*services.TemplateRenderer
	}
)

func init() {
	Register(new(Pages))
}

func (h *Pages) Init(c *services.Container) error {
	h.TemplateRenderer = c.TemplateRenderer
	return nil
}

func (h *Pages) Routes(g *echo.Group) {
	g.GET("/", h.Home).Name = routeNameHome
	g.GET("/about", h.About).Name = routeNameAbout
}

func (h *Pages) Home(ctx echo.Context) error {
	p := page.New(ctx)
	p.Metatags.Description = "Welcome to the homepage."
	p.Metatags.Keywords = []string{"Go", "MVC", "Web", "Software"}
	p.Pager = page.NewPager(ctx, 4)
	p.TemplComponent = pages.Home(h.fetchPosts(&p.Pager))

	return h.RenderPage(ctx, p)
}

// fetchPosts is an mock example of fetching posts to illustrate how paging works
func (h *Pages) fetchPosts(pager *page.Pager) []helpers.Post {
	pager.SetItems(20)
	posts := make([]helpers.Post, 20)

	for k := range posts {
		posts[k] = helpers.Post{
			Title: fmt.Sprintf("Post example #%d", k+1),
			Body:  fmt.Sprintf("Lorem ipsum example #%d ddolor sit amet, consectetur adipiscing elit. Nam elementum vulputate tristique.", k+1),
		}
	}
	return posts[pager.GetOffset() : pager.GetOffset()+pager.ItemsPerPage]
}

func (h *Pages) About(ctx echo.Context) error {
	// A simple example of how the Data field can contain anything you want to send to the templates
	// even though you wouldn't normally send markup like this
	aboutData := helpers.AboutData{
		ShowCacheWarning: true,
		FrontendTabs: []helpers.AboutTab{
			{
				Title: "HTMX",
				Body:  `Completes HTML as a hypertext by providing attributes to AJAXify anything and much more. Visit <a href="https://htmx.org/">htmx.org</a> to learn more.`,
			},
			{
				Title: "Alpine.js",
				Body:  `Drop-in, Vue-like functionality written directly in your markup. Visit <a href="https://alpinejs.dev/">alpinejs.dev</a> to learn more.`,
			},
			{
				Title: "Bulma",
				Body:  `Ready-to-use frontend components that you can easily combine to build responsive web interfaces with no JavaScript requirements. Visit <a href="https://bulma.io/">bulma.io</a> to learn more.`,
			},
		},
		BackendTabs: []helpers.AboutTab{
			{
				Title: "Echo",
				Body:  `High performance, extensible, minimalist Go web framework. Visit <a href="https://echo.labstack.com/">echo.labstack.com</a> to learn more.`,
			},
			{
				Title: "Ent",
				Body:  `Simple, yet powerful ORM for modeling and querying data. Visit <a href="https://entgo.io/">entgo.io</a> to learn more.`,
			},
		},
	}

	p := page.New(ctx)
	p.Title = "About"

	// This page will be cached!
	p.Cache.Enabled = true
	p.Cache.Tags = []string{"page_about", "page:list"}
	p.TemplComponent = pages.About(&aboutData)

	return h.RenderPage(ctx, p)
}
