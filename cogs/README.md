# cogs

cogs is a sub package of cogger that contains prebuilt cogs.  They can be used for many functions such as simple workload, special functionality, and executing system processes.

## Installation

The import path for the package is *gopkg.in/cogger/cogger.v1/cogs*.

To install it, run:

    go get gopkg.in/cogger/cogger.v1/cogs

## Usage

### Simple

~~~ go

simple := cogs.Simple(ctx, func()error{
	//do some work
	return nil
})
~~~

### WithTimeout

WithTimeout creates a simple cog with it's own sperate timeout timer.  This timeout can only be used to make the timeout shorter then global not longer.

~~~ go
withtimeout := cogs.WithTimeout(ctx,func()error{return nil},1*time.Second)
~~~

### Retry

Retry creates a cog that will retry the payload if cogs.ErrRetry is returned up to the maxium attempts.

~~~ go
retry := cogs.Retry(ctx, work func() error{
	//do work
	if somethingBad{
		return cogs.ErrRetry
	}
	return nil
}, 10)
~~~

### Return Error

ReturnErr creates a cog that returns the error passed to it.

~~~ go
returnErr := cogs.ReturnErr(errors.New("some error"))
~~~

### NoOp

NoOp creats a cog that does nothing.

~~~ go
noop := cogs.NoOp()
~~~

### Command

Command creates a cog that executes the command defined as the command and args or the raw command passed it.

~~~ go
generic := cogs.Command(ctx, "command", "arg1","arg2")


buf := &bytes.Buffer{}
withoutput := cogs.CommandWithOutput(ctx, buf, "command",[]string{"arg1","arg2"}...)

cmd := exec.Command("command",[]string{"arg1","arg2"})
raw := cogs.ExecuteCommand(ctx,cmd)
~~~

### DeferredCreate

DeferredCreate defers the time of creation of a cog until right before it needs to execute.  This will allow you to set parameters or other variables in prevous cogs in a series without having to do change the workflow into multiple workflows.

~~~ go
ctx := context.Background()

var subWorkers []cogger.Cog
wait.Resolve(ctx, 
	order.series(ctx,
		cogs.Simple(ctx, func()error{
			//get something from db
			//split each chunk into its own worker
			for _, subStuff := range stuffFromDB{
				workload := subStuff
				subWorkers = append(subWorkers,cog.Simple(ctx,func()error{
					//do work on workload
					return nil
				})
			}
		}),
		cogs.DeferredCreate(func()cogger.Cog{
			return order.Parallel(ctx, subWorkers...)
		}),
	),
)
~~~