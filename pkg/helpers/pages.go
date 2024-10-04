package helpers

type Post struct {
	Title string
	Body  string
}

type AboutData struct {
	ShowCacheWarning bool
	FrontendTabs     []AboutTab
	BackendTabs      []AboutTab
}

type AboutTab struct {
	Title string
	Body  string
}
