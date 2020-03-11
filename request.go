package wine

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gopub/gox"
	"github.com/gopub/wine/internal/io"
	"github.com/gopub/wine/internal/path"
	"github.com/gopub/wine/mime"
)

// Request is a wrapper of http.Request, aims to provide more convenient interface
type Request struct {
	request     *http.Request
	params      gox.M
	body        []byte
	contentType string
}

// Request returns original http request
func (r *Request) Request() *http.Request {
	return r.request
}

// Params returns request parameters
func (r *Request) Params() gox.M {
	return r.params
}

// Body returns request body
func (r *Request) Body() []byte {
	return r.body
}

// ContentType returns request's content type
func (r *Request) ContentType() string {
	return r.contentType
}

// Authorization returns request's Authorization in header
func (r *Request) Authorization() string {
	return r.request.Header.Get("Authorization")
}

// Bearer returns bearer token in header
func (r *Request) Bearer() string {
	s := r.Authorization()
	strs := strings.Split(s, " ")
	if len(strs) != 2 {
		return ""
	}
	if strs[0] == "Bearer" {
		return strs[1]
	}
	return ""
}

func (r *Request) BasicAccount() (user string, password string) {
	s := r.Authorization()
	strs := strings.Split(s, " ")
	if len(strs) != 2 {
		return
	}
	if strs[0] != "Basic" {
		return
	}
	b, err := base64.StdEncoding.DecodeString(strs[1])
	if err != nil {
		logger.Errorf("Decode base64 string %s: %v", strs[1], err)
		return
	}
	userAndPass := strings.Split(string(b), ":")
	if len(userAndPass) != 2 {
		return
	}
	return userAndPass[0], userAndPass[1]
}

func (r *Request) NormalizedPath() string {
	return path.NormalizeRequestPath(r.request)
}

func parseRequest(r *http.Request, maxMem gox.ByteUnit) (*Request, error) {
	params, body, err := io.ReadRequest(r, maxMem)
	if err != nil {
		return nil, fmt.Errorf("read request: %w", err)
	}
	return &Request{
		request:     r,
		params:      params,
		body:        body,
		contentType: mime.GetContentType(r.Header),
	}, nil
}
