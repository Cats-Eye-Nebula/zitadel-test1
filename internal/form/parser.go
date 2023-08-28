package form

import (
	"net/http"
	"net/url"

	"github.com/zitadel/zitadel/internal/errors"

	"github.com/gorilla/schema"
)

type Parser struct {
	decoder *schema.Decoder
}

func NewParser() *Parser {
	d := schema.NewDecoder()
	d.IgnoreUnknownKeys(true)
	return &Parser{d}
}

func (p *Parser) Parse(r *http.Request, data interface{}) error {
	_, err := p.ParseWithFormData(r, data)
	return err
}

func (p *Parser) ParseWithFormData(r *http.Request, data interface{}) (url.Values, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, errors.ThrowInternal(err, "FORM-lCC9zI", "Errors.Internal")
	}

	return r.Form, p.decoder.Decode(data, r.Form)
}
