//package ifilter filters a collection of unknown values by a known interface
package ifilter

import (
	"errors"
	"fmt"
	"reflect"
)

// Collection is a slice of unknown interfaces to be processed.
// It is just an alias of []interface{}, so you can add/remove
// element to it like you would do to a normal slice.
type Collection []interface{}

// Option extends and/or modifies the behavior of Filter.
// Currently no options are available.
type Option interface {
	unimplemented()
}

// Filter examines every element in the collection and callback
// with elements that implement the desired interface. The callback
// argument must be a valid function with the given signature:
//
// func(i []I) error
//
// "I" must be a valid interface, such as io.reader or error.
//
// Callback can optionally return an error. If so, the error will be relayed to the return value of Filter,
// along with any invalid arguments errors raised while filtering.
//
//
func (c Collection) FilterSlice(callback interface{}, Opts ...Option) error {
	cType := reflect.TypeOf(callback)
	if cType == nil {
		return errors.New("callback must not be an untyped nil")
	}
	if cType.Kind() != reflect.Func {
		return fmt.Errorf("callback must be a function, got %v (type %v)", callback, cType)
	}
	if cType.NumIn() != 1 {
		return fmt.Errorf("callback must have exactly 1 argument, got %d", cType.NumIn())
	}
	argType := cType.In(0)
	if argType.Kind() != reflect.Slice {
		return errors.New("argument in callback must be slice")
	}
	if argType.Elem().Kind() != reflect.Interface {
		return errors.New("argument in callback must be a slice of interface")
	}
	iType := argType.Elem()

	var params []reflect.Value
	for i := range c {
		if !implements(c[i], iType) {
			continue
		}
		params = append(params, reflect.ValueOf(c[i]))
	}
	sValue := reflect.MakeSlice(argType, 0, len(params))
	sValue = reflect.Append(sValue, params...)
	outVals := reflect.ValueOf(callback).Call([]reflect.Value{sValue})
	if len(outVals) == 0 {
		return nil
	}
	if last := outVals[len(outVals)-1]; isError(last.Type()) {
		if err, _ := last.Interface().(error); err != nil {
			return err
		}
	}
	return nil
}

// Filter examines every element in the collection and callback
// with elements that implement the desired interface. The callback
// argument must be a valid function with the given signature:
//
// func(i I) error
//
// "I" must be a valid interface, such as io.reader or error.
//
// The callback will be fired for every element implementing I.
// The callback can optionally return an error. If so, the error
// will be relayed to the return value of Filter,
// along with any invalid arguments errors raised while filtering.
//
func (c Collection) Filter(callback interface{}, Opts ...Option) error {
	cType := reflect.TypeOf(callback)
	if cType == nil {
		return errors.New("callback must not be an untyped nil")
	}
	if cType.Kind() != reflect.Func {
		return fmt.Errorf("callback must be a function, got %v (type %v)", callback, cType)
	}
	if cType.NumIn() != 1 {
		return fmt.Errorf("callback must have exactly 1 argument, got %d", cType.NumIn())
	}
	argType := cType.In(0)
	if argType.Kind() != reflect.Interface {
		return errors.New("argument in callback must be a slice of interface")
	}

	for i := range c {
		if !implements(c[i], argType) {
			continue
		}
		outVals := reflect.ValueOf(callback).Call([]reflect.Value{reflect.ValueOf(c[i])})
		if len(outVals) == 0 {
			continue
		}
		if last := outVals[len(outVals)-1]; isError(last.Type()) {
			if err, _ := last.Interface().(error); err != nil {
				return err
			}
		}
	}
	return nil
}

func isError(t reflect.Type) bool {
	return t.Implements(reflect.TypeOf((*error)(nil)).Elem())
}

func implements(i interface{}, iType reflect.Type) bool {
	t := reflect.TypeOf(i)
	if t == nil {
		return false
	}
	return t.Implements(iType)
}
