package testing

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Arrange(name string, u []int) error {
	testData := "data"

	output := fmt.Sprintf("%s/%s.dat", testData, name)

	lines := make([]string, 0)

	for i, file := range u {
		f, err := os.Open(fmt.Sprintf("%s/%s-%d.csv", testData, name, file))
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		line := 0
		for scanner.Scan() {
			l := scanner.Text()
			p := strings.Split(l, ",")
			if len(p) != 3 {
				return fmt.Errorf("invalid line: %s", l)
			}
			if len(lines) <= line {
				prefix := ""
				for j := 0; j < i; j++ {
					prefix += "\t"
				}
				lines = append(lines, prefix+p[2])
			} else {
				lines[line] += "\t" + p[2]
			}
			line++
		}

		if err := scanner.Err(); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()

	return nil
}
