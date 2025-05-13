package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
)

type Method struct {
	name         string
	fn           reflect.Value
	argTypes     []reflect.Type
	returnTypes  []reflect.Type
	returnsError bool
}

func NewMethod(name string, fn any) (*Method, error) {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("is not a function")
	}

	if fnType.NumOut() < 1 {
		return nil, fmt.Errorf("must return at least one value")
	}

	returnsError := false
	for i := range fnType.NumOut() {
		if fnType.Out(i).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			if i != fnType.NumOut()-1 {
				return nil, fmt.Errorf("the error must be the last return value")
			}

			returnsError = true
			break
		}
	}

	method := &Method{
		name:         name,
		fn:           fnValue,
		argTypes:     make([]reflect.Type, fnType.NumIn()),
		returnTypes:  make([]reflect.Type, fnType.NumOut()),
		returnsError: returnsError,
	}

	for i := range fnType.NumIn() {
		method.argTypes[i] = fnType.In(i)
	}

	for i := range fnType.NumOut() {
		method.returnTypes[i] = fnType.Out(i)
	}

	return method, nil
}

func (m *Method) Call(args ...any) (results []any, err error) {
	// Catch panic while running the callback.
	defer func() {
		if _err := recover(); _err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Println("method " + m.name + " crashed: " + fmt.Sprintf("%v\n%s", err, buf))
			err = fmt.Errorf("method handler crashed")
		}
	}()

	if len(args) != len(m.argTypes) {
		return nil, fmt.Errorf("incorrect number of arguments, expected %d", len(m.argTypes))
	}

	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		t := reflect.TypeOf(arg)
		if !t.ConvertibleTo(m.argTypes[i]) {
			return nil, fmt.Errorf("incorrect type for argument %d, expected %s, got %s", i, m.argTypes[i], t)
		}
		in[i] = reflect.ValueOf(arg).Convert(m.argTypes[i])
	}

	out := m.fn.Call(in)

	if m.returnsError {
		if !out[len(out)-1].IsNil() {
			return nil, out[len(out)-1].Interface().(error)
		}
	}

	l := len(out)
	if m.returnsError {
		l--
	}

	results = make([]any, l)
	for i := range results {
		results[i] = out[i].Interface()
	}

	return results, nil
}

func (m *Method) UnmarshalArgs(raw *json.RawMessage) ([]any, error) {
	if len(m.argTypes) == 0 {
		return nil, nil
	}

	args := make([]any, len(m.argTypes))
	if err := json.Unmarshal(*raw, &args); err != nil {
		return nil, err
	}

	return args, nil
}

func (m *Method) MarshalResults(results []reflect.Value) ([]byte, error) {
	if len(results) == 0 {
		return nil, nil
	}

	rawJSON, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return rawJSON, nil
}
