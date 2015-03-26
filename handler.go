package cogger

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

//Handler defines a interface that extends net.Handler interface to have the ability to modify the context of a function
type Handler interface {
	SetTimeout(time.Duration) Handler
	AddContext(...func(context.Context, *http.Request) context.Context) Handler
	SetHandler(func(context.Context, http.ResponseWriter, *http.Request) int) Handler
	ServeHTTP(http.ResponseWriter, *http.Request)
	Test(context.Context, http.ResponseWriter, *http.Request) int
}

type defaultHandler struct {
	f           func(context.Context, http.ResponseWriter, *http.Request) int
	timeout     time.Duration
	ctxCreators []func(context.Context, *http.Request) context.Context
}

func noOp(ctx context.Context, writer http.ResponseWriter, req *http.Request) int {
	return http.StatusInternalServerError
}

//WithHandler is a alias for cogger.NewHandler().SetHandler(f)
func WithHandler(f func(context.Context, http.ResponseWriter, *http.Request) int) Handler {
	h := NewHandler()
	h.SetHandler(f)
	return h
}

//NewHandler creats a Handler interface with the default implementation
func NewHandler() Handler {
	return &defaultHandler{
		f:           noOp,
		timeout:     30 * time.Second,
		ctxCreators: []func(context.Context, *http.Request) context.Context{},
	}
}

func (h *defaultHandler) SetTimeout(t time.Duration) Handler {
	h.timeout = t
	return h
}

func (h *defaultHandler) AddContext(ctxCreators ...func(context.Context, *http.Request) context.Context) Handler {
	h.ctxCreators = append(h.ctxCreators, ctxCreators...)
	return h
}

func (h *defaultHandler) SetHandler(f func(context.Context, http.ResponseWriter, *http.Request) int) Handler {
	h.f = f
	return h
}

func (h *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(base(r), h.timeout)
	defer cancel()
	for _, creator := range h.ctxCreators {
		ctx = creator(ctx, r)
	}
	status := h.f(ctx, w, r)
	if status >= http.StatusBadRequest {
		http.Error(w, http.StatusText(status), status)
	}
}

func (h *defaultHandler) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) int {
	return h.f(ctx, w, r)
}
