package metrics

import (
	"bufio"
	"os"
	"testing"
)

func TestMetrics(t *testing.T) {
	f, err := os.Open("./metrics_output.txt")
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

	output, err := os.OpenFile("./output.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("open output file failed: %v", err)
	}
	defer output.Close()
	for scanner.Scan() {
		line := scanner.Text()
		res, err := rule.Match(line)
		if err != nil {
			continue
		} else {
			t.Logf("matched: %s", res.DumpString())
			// _, _ = output.WriteString(res.DumpString() + "\n")
		}
	}
	var res ResourceList = ConvertToResource(rule)
	_, _ = output.WriteString(res.JSON() + "\n")
	t.Logf(res.JSON())
}
