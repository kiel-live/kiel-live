package testing

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func Arrange(t *testing.T, name string) {
	testData := "data/test-1/"

	files := []string{"1.csv", "10.csv", "100.csv", "1000.csv"}
	output := fmt.Sprintf("data/test-1/%s.dat", name)

	lines := make([]string, 0)

	for i, file := range files {
		f, err := os.Open(testData + name + "-" + file)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		line := 0
		for scanner.Scan() {
			l := scanner.Text()
			p := strings.Split(l, ",")
			if len(p) != 3 {
				t.Fatalf("invalid line %d in %s", line, file)
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
			t.Fatal(err)
		}
	}

	f, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()

	fmt.Println("Done")
}
