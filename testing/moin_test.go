package testing

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func ppTime(name string) {
	fmt.Printf("%s: %d\n", name, time.Now().UnixMicro())
}

func TestMoin(t *testing.T) {
	var poc Poc
	// poc = &graph{}
	poc = &jsonrpc{}

	testSet := &TestSet{
		ID:        "last",
		Latitude:  54.31981897337084,
		Longitude: 10.182968719044112,
	}

	connectingWG := &sync.WaitGroup{}
	connectingWG.Add(1)
	doneWG := &sync.WaitGroup{}
	doneWG.Add(1)
	go func() {
		err := poc.WaitForMessage(nil, connectingWG, func(_ string) {
			timeDiff := time.Since(testSet.StartTime)
			fmt.Printf("received: %f\n", float64(timeDiff.Microseconds())/1000.0)
			fmt.Printf("%s: %d\n", "received", time.Now().UnixMicro())
		})
		require.NoError(t, err)
		doneWG.Done()
	}()

	go func() {
		connectingWG.Wait()
		err := poc.SendData(testSet)
		require.NoError(t, err)
	}()

	doneWG.Wait()

	fmt.Println("done")
}
