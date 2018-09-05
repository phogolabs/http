package httpr

import (
	"net/http"
	"reflect"

	"github.com/go-chi/render"
	"github.com/go-playground/form"
	validator "gopkg.in/go-playground/validator.v9"
)

// DefaultValidator is the default payload validator
var DefaultValidator = validator.New()

func init() {
	render.Decode = DefaultDecoder
}

// DefaultDecoder is the default decoder
func DefaultDecoder(r *http.Request, v interface{}) error {
	var err error

	switch render.GetRequestContentType(r) {
	case render.ContentTypeForm:
		decoder := form.NewDecoder()

		if err = r.ParseForm(); err != nil {
			return err
		}

		return decoder.Decode(v, r.Form)
	default:
		err = render.DefaultDecoder(r, v)
	}

	return err
}

// Bind binds a request into a struct
func Bind(r *http.Request, binder render.Binder) error {
	if err := render.Bind(r, binder); err != nil {
		return err
	}

	return Validate(r, binder)
}

// Validate validates a data
func Validate(r *http.Request, data interface{}) error {
	validator := &(*DefaultValidator)
	validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		switch render.GetRequestContentType(r) {
		case render.ContentTypeJSON:
			return field.Tag.Get("json")
		case render.ContentTypeXML:
			return field.Tag.Get("xml")
		default:
			return field.Name
		}
	})

	if err := validator.StructCtx(r.Context(), data); err != nil {
		return err
	}

	return nil
}
