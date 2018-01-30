# wait

Wait is a sub package of cogger that helps programs wait for completion of cogger workloads. 

## Installation

The import path for the package is *gopkg.in/cogger/cogger.v1/wait*.

To install it, run:

    go get gopkg.in/cogger/cogger.v1/wait

## Usage
### Resolve Single

Single resolvers are concerned with executing the returning errors.  They should be used to execute long chains of cogs.

#### Resolve

Resolve waits until all errors are returned from the cog and then returns them as an array of errors.

~~~ go
var errs []error
errs = wait.Resolve(ctx, cog1)
~~~

#### NoBlock

NoBlock executes the cog and instantly returns control to the initial function.

~~~ go
wait.NoBlock(ctx, cog1)
~~~

#### Completed

Completed executes the cog and instantly returns control to the initial function but calls a call back when it is compelted.

~~~ go
wait.Completed(ctx, cog1, func(){//call me when it completes})
~~~

### Resolve Multiple

Multiple resolvers do not actually resolve cogs but determine what a resolution between multiple cogs looks like.  They should be called in conjuction with a single resolver.

#### All

All will execute all cogs provided in parallel and will return when the first one fails or all succeed

~~~ go

allCog := wait.All(ctx, cog1, cog2, cog3)

~~~

#### Any

Any returns the first cog to finish successfully or ErrNonePassed when all fail

~~~ go

anyCog := wait.Any(ctx, cog1, cog2, cog3)
~~~

#### Settle

Settle will execute all cogs in parallel and return in cog order all the states of the code when all are finished.

~~~ go

cogs := wait.Settle(ctx, cog1, cog2, cog3)
errs := wait.Resolve(ctx, cogs)
//errs[0] cog1 error or nil
//errs[1] cog2 error or nil
//errs[2] cog3 error or nil

~~~
