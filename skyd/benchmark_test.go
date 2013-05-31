package skyd

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var BENCHMARK_PATH = "/tmp/skydb-benchmark-data"
const benchmarkEventCount = 5000000
const benchmarkBatchSize = 1000

func TestGenerateData(t *testing.T) {
	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	
	t0 := time.Now()
	warn("benchmarking...")
	warn("%s", path)
	withBenchmarkData(path, benchmarkEventCount, benchmarkBatchSize, func(s *Server) {
		res := benchmarkEventsCount(s)
		if res != benchmarkEventCount {
			t.Errorf("withBenchmarkData did not generate %d events before running block, returned %v", benchmarkEventCount, res)
		}
		warn("benchmark done")
		warn("total:      %d events", res)
		warn("throughput: %d usec/event", (time.Now().Sub(t0).Nanoseconds() / int64(res)) / 1000)
		warn("")
	})
}

