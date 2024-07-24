package testing

import (
	"sync"
)

type jsonrpc struct{}

func (g *jsonrpc) Name() string {
	return "jsonrpc"
}

func (g *jsonrpc) SendData(testSet *TestSet) error {
	return nil
}

func (g *jsonrpc) WaitForMessage(testSets []*TestSet, connectingWG *sync.WaitGroup, i int, id int, write func(s string)) error {
	return nil
}
