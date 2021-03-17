package form

import (
	"net/url"

	"github.com/go-playground/form/v4"
)

var formDecoder = form.NewDecoder()

func Decode(v interface{}, values url.Values) error {
	return formDecoder.Decode(v, values)
}
