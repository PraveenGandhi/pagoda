package funcs

import (
	"fmt"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mikestefanello/pagoda/config"

	"github.com/stretchr/testify/assert"
)

func TestHasField(t *testing.T) {
	type example struct {
		name string
	}
	var e example
	f := new(Funcs)
	assert.True(t, f.HasField(e, "name"))
	assert.False(t, f.HasField(e, "abcd"))
}

func TestLink(t *testing.T) {
	f := new(Funcs)

	link := string(f.Link("/abc", "Text", "/abc"))
	expected := `<a class="is-active" href="/abc">Text</a>`
	assert.Equal(t, expected, link)

	link = string(f.Link("/abc", "Text", "/abc", "first", "second"))
	expected = `<a class="first second is-active" href="/abc">Text</a>`
	assert.Equal(t, expected, link)

	link = string(f.Link("/abc", "Text", "/def"))
	expected = `<a class="" href="/abc">Text</a>`
	assert.Equal(t, expected, link)
}

func TestFile(t *testing.T) {
	f := new(Funcs)

	file := f.File("test.png")
	expected := fmt.Sprintf("/%s/test.png?v=%s", config.StaticPrefix, CacheBuster)
	assert.Equal(t, expected, file)
}

func TestUrl(t *testing.T) {
	f := new(Funcs)
	f.web = echo.New()
	f.web.GET("/mypath/:id", func(c echo.Context) error {
		return nil
	}).Name = "test"
	out := f.URL("test", 5)
	assert.Equal(t, "/mypath/5", out)
}
