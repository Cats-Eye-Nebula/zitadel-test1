package saml

import (
	"bytes"
	"context"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Session = (*Session)(nil)

// Session is the [idp.Session] implementation for the SAML provider.
type Session struct {
	ServiceProvider *samlsp.Middleware
	state           string

	RequestID string
	Request   *http.Request
}

// GetAuthURL implements the [idp.Session] interface.
func (s *Session) GetAuth(ctx context.Context) (string, bool) {
	url, _ := url.Parse(s.state)
	resp := NewTempResponseWriter()

	request := &http.Request{
		URL: url,
	}
	s.ServiceProvider.HandleStartAuthFlow(
		resp,
		request.WithContext(ctx),
	)

	if location := resp.Header().Get("Location"); location != "" {
		return idp.Redirect(location)
	}
	return idp.Form(string(resp.content.Bytes()))
}

// FetchUser implements the [idp.Session] interface.
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	if s.RequestID == "" || s.Request == nil {
		return nil, errors.ThrowInvalidArgument(nil, "SAML-d09hy0wkex", "Errors.Intent.ResponseInvalid")
	}

	assertion, err := s.ServiceProvider.ServiceProvider.ParseResponse(s.Request, []string{s.RequestID})
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "SAML-nuo0vphhh9", "Errors.Intent.ResponseInvalid")
	}

	userMapper := NewUser()
	userMapper.SetID(assertion.Subject.NameID)
	for _, statement := range assertion.AttributeStatements {
		for _, attribute := range statement.Attributes {
			values := make([]string, len(attribute.Values))
			for i := range attribute.Values {
				values[i] = attribute.Values[i].Value
			}
			userMapper.Attributes[attribute.Name] = values
		}
	}
	return userMapper, nil
}

type TempResponseWriter struct {
	header  http.Header
	content *bytes.Buffer
}

func (w *TempResponseWriter) Header() http.Header {
	return w.header
}

func (w *TempResponseWriter) Write(content []byte) (int, error) {
	return w.content.Write(content)
}

func (w *TempResponseWriter) WriteHeader(statusCode int) {}

func NewTempResponseWriter() *TempResponseWriter {
	return &TempResponseWriter{
		header:  map[string][]string{},
		content: bytes.NewBuffer([]byte{}),
	}
}
