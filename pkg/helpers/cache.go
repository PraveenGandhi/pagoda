package helpers

import "github.com/mikestefanello/pagoda/pkg/form"

type CacheForm struct {
	Value string `form:"value"`
	form.Submission
}
