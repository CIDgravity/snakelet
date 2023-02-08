// handy validators file
// see example usage here https://github.com/go-playground/validator/blob/master/_examples/custom/main.go
package snakelet

import (
	"net/url"
	"reflect"
)

func ValidateUrl(field reflect.Value) interface{} {
	url, err := url.Parse(field.String())
	if err != nil {
		return err
	} else {
		return url
	}
}
