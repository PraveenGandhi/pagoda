package helpers

import "github.com/mikestefanello/pagoda/pkg/form"

type TaskForm struct {
	Delay   int    `form:"delay" validate:"gte=0"`
	Message string `form:"message" validate:"required"`
	form.Submission
}
