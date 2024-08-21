package usage

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type Usage struct {
	name     string
	done     context.CancelFunc
	waitDone chan struct{}
	f        *os.File
	w        *bufio.Writer
}

func NewUsage(name string) *Usage {
	u := &Usage{
		name: name,
	}

	filePath := fmt.Sprintf("./testing/data/usage-%s.csv", name)
	os.Remove(filePath)

	var err error
	u.f, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}

	u.w = bufio.NewWriter(u.f)

	return u
}

func (u *Usage) getUsage() (string, error) {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return "", err
	}

	cpu, err := p.CPUPercent()
	if err != nil {
		return "", err
	}

	mem, err := p.MemoryInfo()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d,%f\n", mem.RSS, cpu), nil
}

func (u *Usage) Collect(amountClients string) error {
	if u.done != nil {
		u.done()
		<-u.waitDone
		fmt.Printf("amount: %s\n", amountClients)
		_, err := u.w.WriteString("\n\n") // 2 empty lines for next test
		if err != nil {
			return err
		}
	}

	if amountClients == "end" {
		u.w.Flush()
		u.f.Close()
		return nil
	}

	ctx, done := context.WithCancel(context.Background())
	u.done = done
	u.waitDone = make(chan struct{})

	d := 1 * time.Microsecond

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	start := time.Now()

out:
	for {
		select {
		case <-ctx.Done():
			break out

		case <-ticker.C:
			s, err := u.getUsage()
			if err != nil {
				return err
			}
			diff := time.Since(start)
			_, err = u.w.WriteString(fmt.Sprintf("%f,%s,%s", float64(diff.Microseconds()/1000.0), amountClients, s))
			if err != nil {
				return err
			}
		}
	}

	err := u.w.Flush()
	if err != nil {
		return err
	}

	close(u.waitDone)

	return nil
}

func perf(name string) *os.File {
	f, err := os.Create(fmt.Sprintf("data/cpu-%s.prof", name))
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}

	return f
}

func perfEnd(name string) {
	defer pprof.StopCPUProfile()

	f1, err := os.Create(fmt.Sprintf("data/heap-%s.prof", name))
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f1.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f1); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
}
