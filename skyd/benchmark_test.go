package skyd

import (
	"io/ioutil"
	"os"
	"testing"
)

var BENCHMARK_PATH = "/tmp/skydb-benchmark-data"

func TestGenerateData(t *testing.T) {

	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)

	withBenchmarkData(path, 10, func(s *Server) {
		res := benchmarkEventsCount(s)
		if res != 10 {
			t.Errorf("withBenchmarkData did not generate 10 events before running block, returned %v", res)
		}
	})

}
