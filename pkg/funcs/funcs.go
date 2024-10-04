package funcs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/mikestefanello/pagoda/config"
)

var (
	// CacheBuster stores a random string used as a cache buster for static files.
	CacheBuster = random.String(10)
)

type Funcs struct {
	web *echo.Echo
}

// NewFuncMap provides a template function map
func NewFuncs(web *echo.Echo) *Funcs {
	return &Funcs{web: web}
}

// HasField checks if an interface contains a given field
func (fm *Funcs) HasField(v any, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

// File appends a cache buster to a given filepath so it can remain cached until the app is restarted
func (fm *Funcs) File(filepath string) string {
	return fmt.Sprintf("/%s/%s?v=%s", config.StaticPrefix, filepath, CacheBuster)
}

// Link outputs HTML for a link element, providing the ability to dynamically set the active class
func (fm *Funcs) Link(routeName, text, currentPath string, classes ...string) string {
	url := fm.URL(routeName)
	if currentPath == string(url) {
		classes = append(classes, "is-active")
	}

	return fmt.Sprintf(`<a class="%s" href="%s">%s</a>`, strings.Join(classes, " "), url, text)
}

// URL generates a URL from a given route name and optional parameters
func (fm *Funcs) URL(routeName string, params ...any) templ.SafeURL {
	return templ.URL(fm.web.Reverse(routeName, params...))
}
