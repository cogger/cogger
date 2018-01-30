package cogger

import (
	"net/http"

	"golang.org/x/net/context"
)

var base = func(request *http.Request) context.Context {
	return context.Background()
}
