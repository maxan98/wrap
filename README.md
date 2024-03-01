# Wrap
## Wrap is a simple golang utility used to wrap (decorate) just any function and add retry mechanism

## !NOT! for production, just 4 fun

### Based on github.com/avast/retry-go

Small lib that allows you to apply retries to literally any function, regardless of return values and input params.
This is not an exact decorator as it modifies return type.

You can pass in any function with any input params(types/amount) and return values(types/amount).
As a result of a decorated function execution you will receive an array containing all original return values ([]reflect.Values) and err (wrapped errors from all retries)

Retry options can be passed as a last argument in Retry function or modified globally.

#### Usage

`go get github.com/maxan98/wrap`

Example usage:

```go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/maxan98/wrap"
)

var errRetry = fmt.Errorf("retry")

type Example struct{}

func (v Example) Test(val string) (string, error) {
	fmt.Println(time.Now())
	return val, errRetry
}

func main() {
	wrap.DefaultRetryIf = func(err error) bool {
		return errors.Is(err, errRetry)
	}
	wrap.DefaultAttempts = 3

	a := Example{}
	wrapped, err := wrap.Retry(a.Test, []any{"sdf"})
	fmt.Println(err) // returns error if unable to wrap, e.g. 1st arg is not a func or len(2nd arg) != len(actual func's args)
	res, err := wrapped()

	fmt.Println(res)
	fmt.Println("err", err)
}
```