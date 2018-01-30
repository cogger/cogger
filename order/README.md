# order

order is a sub package of cogger that determines how multipel cogs would run together. 

## Installation

The import path for the package is *gopkg.in/cogger/cogger.v1/order*.

To install it, run:

    go get gopkg.in/cogger/cogger.v1/order

## Usage

### Series

Series runs cogs in series.

~~~ go

series := order.Series(ctx,
	cog1,
	cog2,
	cog3,
)

~~~

### Parallel

Parallel runs cogs in parallel.

~~~ go
parallel := order.Parellel(ctx,
	cog1,
	cog2,
	cog3,
)
~~~

### If

If only runs the cog if the test passes.

~~~ go

ifcog := order.If(ctx,
	func(ctx context.Context)bool{return true},
	cog1,
)

~~~

## Complex Usage

You can stack orders in any pattern you desire.

~~~ go

workload := order.Series(ctx,
	cog_1,
	order.Parallel(ctx, 
		order.Series(ctx,
			cog_2_1_1,
			cog_2_1_2,
		),
		cog_2_2,
		order.If(ctx,
			func(ctx context.Context)bool{return somePreviousWorkDidsomething},
			cog_2_3
		),
	),
	cog_3,
)

~~~