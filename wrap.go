package wrap

import (
	"fmt"
	"reflect"

	"github.com/avast/retry-go"
)

type Response struct {
	Values []reflect.Value
	Types  []reflect.Type
}

var (
	DefaultAttempts      = retry.DefaultAttempts
	DefaultDelay         = retry.DefaultDelay
	DefaultMaxJitter     = retry.DefaultMaxJitter
	DefaultOnRetry       = retry.DefaultOnRetry
	DefaultRetryIf       = retry.DefaultRetryIf
	DefaultDelayType     = retry.DefaultDelay
	DefaultLastErrorOnly = retry.DefaultLastErrorOnly
	DefaultContext       = retry.DefaultContext
)

// Retry accepts func with any arguments but with only one error return type in return set
// If func return multiple error types - only last one will be seen in wrapped errors set after retries
func Retry(inner any, args []any, opts ...retry.Option) (func() (Response, error), error) {

	retry.DefaultAttempts = DefaultAttempts
	retry.DefaultDelay = DefaultDelay
	retry.DefaultMaxJitter = DefaultMaxJitter
	retry.DefaultOnRetry = DefaultOnRetry
	retry.DefaultRetryIf = DefaultRetryIf
	retry.DefaultDelay = DefaultDelayType
	retry.DefaultLastErrorOnly = DefaultLastErrorOnly
	retry.DefaultContext = DefaultContext

	t := reflect.TypeOf(inner)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("inner is not a func")
	} else if t.NumIn() != len(args) {
		return nil, fmt.Errorf("wrong number of params")
	}
	argsList := make([]reflect.Value, len(args))
	for i, arg := range args {
		argsList[i] = reflect.ValueOf(arg)
	}
	var arr Response
	b := func() (Response, error) {
		err := retry.Do(func() error {
			var result []reflect.Value
			arr = Response{}
			result = reflect.ValueOf(inner).Call(argsList)
			var retErr error

			for _, val := range result {
				if e, ok := val.Interface().(error); ok {
					retErr = e
				} else {
					arr.Values = append(arr.Values, val)
					arr.Types = append(arr.Types, val.Type())
				}
			}
			return retErr
		}, opts...)
		if err != nil {
			return arr, err
		}
		return arr, nil
	}
	return b, nil
}
