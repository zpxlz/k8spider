package metrics

import (
	"bufio"
	"os"
	"testing"
)

func TestMetrics(t *testing.T) {
	f, err := os.Open("./metrics")
	if err != nil {
		t.Fatalf("open file failed: %v", err)
		t.Fail()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	rule := DefaultMatchRules()
	if err := rule.Compile(); err != nil {
		t.Fatalf("compile rule failed: %v", err)
		t.Fail()
	}

	for scanner.Scan() {
		line := scanner.Text()
		res, err := rule.Match(line)
		if err != nil {
			continue
		} else {
			t.Logf("matched: %s", res.DumpString())
		}
	}
}
