package skyd

import (
	"encoding/json"
	"testing"
)

// Ensure that a cluster can serialize its contents.
func TestClusterSerialize(t *testing.T) {
	cluster := NewCluster()
	cluster.addNodeGroup()
	cluster.addNodeGroup()
	cluster.addNode("foo", 0)
	cluster.addNode("bar", 0)
	cluster.addNode("baz", 1)
	s, _ := json.Marshal(cluster.serialize())
	if string(s) != `{"groups":[{"nodes":[{"name":"foo"},{"name":"bar"}],"shards":""},{"nodes":[{"name":"baz"}],"shards":""}]}` {
		t.Fatalf("Unexpected cluster serialization: %v", string(s))
	}
}
