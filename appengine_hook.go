// +build appengine,appenginevm

package cogger

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func init() {
	base = func(request http.Request) context.Context {
		return appengine.NewContext(request)
	}
}
