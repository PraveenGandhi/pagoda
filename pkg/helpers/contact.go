package helpers

import "github.com/mikestefanello/pagoda/pkg/form"

type ContactForm struct {
	Email      string `form:"email" validate:"required,email"`
	Department string `form:"department" validate:"required,oneof=sales marketing hr"`
	Message    string `form:"message" validate:"required"`
	form.Submission
}
