package metrics

import (
	"bufio"
	"os"
	"testing"
)

var rule = []*MetricMatcher{
	NewMetricMatcher("configmap").AddLabel("namespace").AddLabel("configmap"),
	NewMetricMatcher("pod").AddLabel("namespace").AddLabel("pod").AddLabel("host_ip").AddLabel("pod_ip"),
}

func TestMetrics(t *testing.T) {
	f, err := os.Open("./metrics_output.txt")
	if err != nil {
		t.Fatalf("open file failed: %v", err)
		t.Fail()
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for i := range rule {
		err := rule[i].Compile()
		if err != nil {
			t.Fatalf("compile rule failed: %v", err)
			t.Fail()
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		for _, r := range rule {
			res, e := r.Match(line)
			_ = res
			if e != nil {
				continue
			} else {
				t.Log(r.DumpString())
			}
		}
	}
}
