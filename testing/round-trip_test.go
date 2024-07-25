package testing

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Poc interface {
	Name() string
	SendData(testSet *TestSet) error
	WaitForMessage(testSets []*TestSet, connectingWG *sync.WaitGroup, done func(s string)) error
}

func TestRoundTrip(t *testing.T) {
	u := []int{1, 10, 100, 1000}
	// repeat := 10
	var poc Poc
	poc = &graph{}
	poc = &jsonrpc{}

	for _, amountClients := range u {
		f, err := os.OpenFile(fmt.Sprintf("./data/test-1/%s-%d.csv", poc.Name(), amountClients), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		m := sync.Mutex{}
		writer := bufio.NewWriter(f)

		testSets := []*TestSet{
			// {"123", 54.31981897337084, 10.182968719044112, time.Now()},
			{"last", 54.31981897337084, 10.182968719044112, time.Now()},
		}

		for i := 0; i < 10000/amountClients; i++ {
			fmt.Printf("# [%d:%d]\n", amountClients, i)

			connectingWG := &sync.WaitGroup{}
			doneWG := &sync.WaitGroup{}

			connect := func(id int) {
				err := poc.WaitForMessage(testSets, connectingWG, func(testSetID string) {
					for _, testSet := range testSets {
						if testSet.ID == testSetID {
							timeDiff := time.Since(testSet.StartTime)
							// fmt.Printf("[%d:%d] %s: %s\n", i, id, data.MapStopUpdated.ID, timeDiff)
							// fmt.Printf("%d,%d,%f\n", i, id, float64(timeDiff.Microseconds())/1000.0)

							s := fmt.Sprintf("%d,%d,%f\n", i, id, float64(timeDiff.Microseconds())/1000.0)

							m.Lock()
							_, err := writer.WriteString(s)
							require.NoError(t, err)
							m.Unlock()
						}
					}
				})
				require.NoError(t, err)
				doneWG.Done()
			}

			for i := 0; i < amountClients; i++ {
				connectingWG.Add(1)
				doneWG.Add(1)
				go connect(i)
			}

			go func() {
				connectingWG.Wait()
				// fmt.Println("all clients connected")
				if poc.Name() == "graphql" {
					time.Sleep(10 * time.Duration(amountClients) * time.Millisecond) // wait some a bit longer to be really sure everyone is listening
				}

				for _, testSet := range testSets {
					err := poc.SendData(testSet)
					require.NoError(t, err)
				}
			}()

			doneWG.Wait()
			// fmt.Println("all clients done")
		}

		writer.Flush()
	}

	err := Arrange(poc.Name(), u)
	require.NoError(t, err)

	fmt.Println("done")
}
