package service_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kiel-live/kiel-live/jsonrpc/rpc/service"
)

type SampleRPC struct{}

func (t *SampleRPC) Add(a int, b int) int {
	return a + b
}

func (t *SampleRPC) Hello(name string) (string, error) {
	return "Hello, " + name, nil
}

type ComplexType struct {
	Name string
}

func (t *SampleRPC) Complex(name string) (*ComplexType, error) {
	return &ComplexType{
		Name: "Hello, " + name,
	}, nil
}

func (t *SampleRPC) Fails() (*ComplexType, error) {
	return nil, fmt.Errorf("failed")
}

func TestMethod(t *testing.T) {
	rpc := &SampleRPC{}
	s, err := service.NewService(rpc)
	assert.NoError(t, err)

	// Hello
	args := json.RawMessage("[\"World\"]")
	results, err := s.Call("Hello", &args)
	assert.NoError(t, err)
	j, err := json.Marshal(results)
	assert.NoError(t, err)
	assert.JSONEq(t, `["Hello, World"]`, string(j))

	// Add
	args = json.RawMessage("[3, 4]")
	results, err = s.Call("Add", &args)
	assert.NoError(t, err)
	j, err = json.Marshal(results)
	assert.NoError(t, err)
	assert.JSONEq(t, `[7]`, string(j))

	// Complex
	args = json.RawMessage("[\"World\"]")
	results, err = s.Call("Complex", &args)
	assert.NoError(t, err)
	j, err = json.Marshal(results)
	assert.NoError(t, err)
	assert.JSONEq(t, `[{"Name":"Hello, World"}]`, string(j))

	// Fails
	results, err = s.Call("Fails", nil)
	assert.Errorf(t, err, "failed")
	assert.Nil(t, results)
}
