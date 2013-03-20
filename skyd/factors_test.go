package skyd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// Ensure that we can create a new table.
func TestFactorization(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	path = fmt.Sprintf("%v/factors", path)

	factors := NewFactors(path)
	defer factors.Close()
	err = factors.Open()
	if err != nil {
		t.Fatalf("Unable to create factors: %v", err)
	}
	
	num, err := factors.Factorize("foo", "bar", "/index.html")
	if err != nil {
		t.Fatalf("Unable to factorize: %v", err)
	}
	if num != 1 {
		t.Fatalf("Wrong factorization: exp: %v, got: %v", 1, num)
	}
}

