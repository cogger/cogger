# cogger [![GoDoc](https://godoc.org/github.com/cogger/cogger?status.png)](http://godoc.org/github.com/cogger/cogger)

cogger is a package that extends golang.org/x/net/context with additional helper functions to make it easy to implement the context pattern as dicussed at blog.golang.org/context.  It allows you to manage mutliple construction and tear down of go coroutines, scopes items on your context per request and generally makes your program to appear sequential while still being highly concurrent.

## Usage

You can use cogger to help manage http handler functions.

~~~ go
// main.go
package main

import (
	"log"
	"net/http"
	"html"
	"fmt"
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

func fooHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) int{
	if somethingBad() {
		return http.StatusInternalServerError
	}

	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	return http.StatusOK
}

func main() {
  	http.Handle("/foo", cogger.NewHandler().AddContext(...something).SetHandler(fooHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

~~~

You can use cogger without http handler functions.

~~~ go
// main.go
package main

import (
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
)

func main() {
  	cogger.With(func(ctx context.Context){
  		//Do something
  	})
}

~~~

You can manage complex task execution.

~~~ go
// main.go
package main

import (
	"github.com/cogger/cogger/cogs"
	"github.com/cogger/cogger/order"
	"github.com/cogger/cogger/wait"
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
	"time"
)

type Doc struct{}
func getDBDoc() (Doc,error){
	//get from DB
	return Doc{},nil
}

type Result struct{}
func makeResult(doc Doc) (Result, error){
	//make a result
	return Result, nil
}

func saveResultToFile(result Result) error{
	//save the result
	return nil
}

func saveResultToDB(result Result) error{
	//save the result
	return nil
}

func main() {
  	cogger.WithTimeout(func(ctx context.Context){
  		var doc Doc
  		getDocWorker := cogs.Simple(ctx, func() error{
  			var err error
  			doc, err = getDBDoc()
  			return err
  		})

  		var result Result
  		makeResultWorker := cogs.WithTimeout(ctx, func()error{
  			var err error
  			result, err = makeResult(doc)
  			return err
  		},5 * time.Second)

  		err := <- order.Series(wait.Any(getDocWorker, getDocWorker, getDocWorker),
  			makeResultWorker,
  			wait.Settle(ctx, work.Retry(ctx, func()error{
	  			return saveResultToFile(result)
	  		}),
	  		work.Retry(ctx, func()error{
	  			return saveResultToDB(result)
	  		}),
  		)).Do()
  		
  		if err != nil {
  			//Handle the error
  		}
  	},30 * time.Second)
}

~~~

You can set limits on how many cogs you want working at any given point.

~~~ go
// main.go
package main

import (
	"github.com/cogger/cogger/cogs"
	"github.com/cogger/cogger/order"
	"github.com/cogger/cogger/wait"
	"github.com/cogger/cogger"
	"golang.org/x/net/context"
	"time"
)

func doSomething() error{
	return nil
}

func main() {
  	cogger.With(func(ctx context.Context){
  		fiveTotal := limiter.ByCount(5)

		howMany := 30
  		workers := make([]cogger.Cog,howMany)
  		for i := 0; i < howMany; i ++{
	  		worker := cogs.Simple(ctx, func() error{
	  			return doSomething()
	  		})
  			worker.SetLimit(fiveTotal)
  			workers[i] = worker
  		}

  		//Even though this has 30 cogs it will only execute 5 at a time
  		<- order.Parallel(ctx, workers...).Do()
  		
  		onePerSecond := limiter.PerSecond(1)
  		for i := 0; i < howMany; i ++{
	  		worker := cogs.Simple(ctx, func() error{
	  			return doSomething()
	  		})
  			worker.SetLimit(onePerSecond)
  			workers[i] = worker
  		}

  		//Even though this has 30 cogs it will only start 1 per second
  		<- order.Parallel(ctx, workers...).Do()

  	})
}

~~~