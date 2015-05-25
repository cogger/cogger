# limiter

limiter is a sub package of cogger that contains prebuilt cog limiters.  Cog limiters allow global control of resourses used by each function.

## Installation

The import path for the package is *gopkg.in/cogger/cogger.v1/limiter*.

To install it, run:

    go get gopkg.in/cogger/cogger.v1/limiter

## Usage

### ByCount

ByCount creates a limiter that will only allow that number of cogs to be running at any given point.

~~~ go
limit := limiter.ByCount(10)
cog := cogs.Simple(ctx, func()error{return nil})
cog.SetLimit(limit)
~~~

### PerSecond

PerSecond creates a limiter that will allow that number of cogs to be started per second.  It does not assure that only x per second are running only taht x per second have been started.

~~~ go
limit := limiter.PerSecond(10)
cog := cogs.Simple(ctx, func()error{return nil})
cog.SetLimit(limit)
~~~