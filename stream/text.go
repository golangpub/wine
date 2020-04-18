package stream

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/golangpub/types"

	"github.com/golangpub/log"
	"github.com/golangpub/wine"
	"github.com/golangpub/wine/api"
	"github.com/golangpub/wine/mime"
)

const textPacketDelimiter = 0x01

type TextReadCloser interface {
	Read() (string, error)
	io.Closer
}

type TextWriteCloser interface {
	Write(s string) error
	io.Closer
}

type textReadCloser struct {
	body  io.ReadCloser
	buf   *bytes.Buffer
	block []byte
	err   error
}

func newTextReadCloser(body io.ReadCloser) *textReadCloser {
	r := new(textReadCloser)
	r.body = body
	r.buf = new(bytes.Buffer)
	r.block = make([]byte, 1024)
	return r
}

func (r *textReadCloser) Read() (string, error) {
	for {
		p, ok := r.readPacket()
		if ok {
			return p, nil
		}
		if r.err != nil {
			return "", r.err
		}
		n, err := r.body.Read(r.block)
		if n > 0 {
			r.buf.Write(r.block[:n])
		}
		if err != nil {
			r.err = err
		}
	}
}

func (r *textReadCloser) Close() error {
	return r.body.Close()
}

func (r *textReadCloser) readPacket() (string, bool) {
	for _, b := range r.buf.Bytes() {
		if b == textPacketDelimiter {
			p, err := r.buf.ReadBytes(textPacketDelimiter)
			if err != nil {
				log.Errorf("Read bytes: %v", err)
				return "", false
			}
			// Exclude the last byte which is packet delimiter
			return string(p[:len(p)-1]), true
		}
	}
	return "", false
}

type textWriteCloser struct {
	w    http.ResponseWriter
	done chan<- interface{}
}

func newTextWriteCloser(w http.ResponseWriter, done chan<- interface{}) *textWriteCloser {
	return &textWriteCloser{
		w:    w,
		done: done,
	}
}

func (w *textWriteCloser) Write(s string) error {
	p := []byte(s)
	p = append(p, textPacketDelimiter)
	_, err := w.w.Write(p)
	if err != nil {
		return fmt.Errorf("write all: %w", err)
	}
	if flusher, ok := w.w.(http.Flusher); ok {
		flusher.Flush()
	}
	if e, ok := w.w.(interface{ Error() error }); ok {
		if e.Error() != nil {
			return e.Error()
		}
	}
	return nil
}

func (w *textWriteCloser) Close() error {
	close(w.done)
	return nil
}

func NewTextReader(client *http.Client, req *http.Request) (TextReadCloser, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		err = api.ParseResult(resp, nil, true)
		if err != nil {
			return nil, fmt.Errorf("parse result: %w", err)
		}
		return nil, types.NewError(resp.StatusCode, "unknown error")
	}
	r := newTextReadCloser(resp.Body)
	greeting, err := r.Read()
	if err != nil {
		r.Close()
		return nil, fmt.Errorf("handshake: %w", err)
	}
	if greeting != Greeting {
		r.Close()
		return nil, fmt.Errorf("expect %s, got %s", Greeting, greeting)
	}
	return r, nil
}

func NewTextHandler(serve func(context.Context, TextWriteCloser)) wine.Handler {
	return wine.HandlerFunc(func(ctx context.Context, req *wine.Request, next wine.Invoker) wine.Responder {
		logger := log.FromContext(ctx)
		logger.Debugf("Start")
		defer logger.Debugf("Closed")
		w := wine.GetResponseWriter(ctx)
		w.Header().Set(mime.ContentType, mime.HTML_UTF8)
		done := make(chan interface{})
		tw := newTextWriteCloser(w, done)
		err := tw.Write(Greeting)
		if err != nil {
			logger.Errorf("Handshake: %v", err)
			return wine.Status(http.StatusOK)
		}
		go serve(ctx, tw)
		<-done
		return wine.Status(http.StatusOK)
	})
}
