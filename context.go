package wine

import (
	"context"
	"github.com/golangpub/types"
	"html/template"
	"net/http"
)

type contextKey int

// Context keys
const (
	ckBasicAuthUser contextKey = iota + 1
	ckHTTPResponseWriter
	ckTemplates
	ckSessionID
	ckRemoteAddr
	ckLocation
	ckAccessToken
	ckUserID
	ckTraceID
	ckUser
	ckDeviceID
)

func GetBasicAuthUser(ctx context.Context) string {
	user, _ := ctx.Value(ckBasicAuthUser).(string)
	return user
}

func withBasicAuthUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, ckBasicAuthUser, user)
}

func GetSessionID(ctx context.Context) string {
	sid, _ := ctx.Value(ckSessionID).(string)
	return sid
}

func withSessionID(ctx context.Context, sid string) context.Context {
	return context.WithValue(ctx, ckSessionID, sid)
}

// GetResponseWriter returns response writer from the context
func GetResponseWriter(ctx context.Context) http.ResponseWriter {
	rw, _ := ctx.Value(ckHTTPResponseWriter).(http.ResponseWriter)
	return rw
}

func withResponseWriter(ctx context.Context, rw http.ResponseWriter) context.Context {
	return context.WithValue(ctx, ckHTTPResponseWriter, rw)
}

// GetTemplates returns templates in context
func GetTemplates(ctx context.Context) []*template.Template {
	v, _ := ctx.Value(ckTemplates).([]*template.Template)
	return v
}

func withTemplate(ctx context.Context, templates []*template.Template) context.Context {
	return context.WithValue(ctx, ckTemplates, templates)
}

func GetUserID(ctx context.Context) int64 {
	id, _ := ctx.Value(ckUserID).(int64)
	return id
}

func WithUserID(ctx context.Context, id int64) context.Context {
	if id == 0 {
		return ctx
	}
	return context.WithValue(ctx, ckUserID, id)
}

func GetUser(ctx context.Context) interface{} {
	return ctx.Value(ckUser)
}

func WithUser(ctx context.Context, u interface{}) context.Context {
	return context.WithValue(ctx, ckUser, u)
}

func GetAccessToken(ctx context.Context) string {
	token, _ := ctx.Value(ckAccessToken).(string)
	return token
}

func WithAccessToken(ctx context.Context, token string) context.Context {
	if len(token) == 0 {
		return ctx
	}
	return context.WithValue(ctx, ckAccessToken, token)
}

func GetRemoteAddr(ctx context.Context) string {
	ip, _ := ctx.Value(ckRemoteAddr).(string)
	return ip
}

func WithRemoteAddr(ctx context.Context, addr string) context.Context {
	if len(addr) == 0 {
		return ctx
	}
	return context.WithValue(ctx, ckRemoteAddr, addr)
}

func GetDeviceID(ctx context.Context) string {
	id, _ := ctx.Value(ckDeviceID).(string)
	return id
}

func WithDeviceID(ctx context.Context, deviceID string) context.Context {
	if len(deviceID) == 0 {
		return ctx
	}
	return context.WithValue(ctx, ckDeviceID, deviceID)
}

func GetTraceID(ctx context.Context) string {
	id, _ := ctx.Value(ckTraceID).(string)
	return id
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	if len(traceID) == 0 {
		return ctx
	}
	return context.WithValue(ctx, ckTraceID, traceID)
}

func GetLocation(ctx context.Context) *types.Point {
	id, _ := ctx.Value(ckLocation).(*types.Point)
	return id
}

func WithLocation(ctx context.Context, location *types.Point) context.Context {
	if location == nil {
		return ctx
	}
	return context.WithValue(ctx, ckLocation, location)
}

func DetachContext(ctx context.Context) context.Context {
	newCtx := context.Background()
	if t := GetTemplates(ctx); len(t) != 0 {
		newCtx = withTemplate(ctx, t)
	}
	if u := GetBasicAuthUser(ctx); u != "" {
		newCtx = withBasicAuthUser(ctx, u)
	}
	if sid := GetSessionID(ctx); sid != "" {
		newCtx = withSessionID(ctx, sid)
	}
	if token := GetAccessToken(ctx); token != "" {
		newCtx = WithAccessToken(newCtx, token)
	}
	if deviceID := GetDeviceID(ctx); deviceID != "" {
		newCtx = WithDeviceID(newCtx, deviceID)
	}
	if c := GetLocation(ctx); c != nil {
		newCtx = WithLocation(newCtx, c)
	}
	if addr := GetRemoteAddr(ctx); addr != "" {
		newCtx = WithRemoteAddr(newCtx, addr)
	}
	if traceID := GetTraceID(ctx); traceID != "" {
		newCtx = WithTraceID(newCtx, traceID)
	}
	if uid := GetUserID(ctx); uid > 0 {
		newCtx = WithUserID(newCtx, uid)
	}
	if u := GetUser(ctx); u != nil {
		newCtx = WithUser(newCtx, u)
	}
	return newCtx
}
