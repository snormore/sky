package skyd

import (
	"io/ioutil"
	"os"
	"testing"
)

// Ensure that we can open and close a tablet.
func TestOpen(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	
	tablet := NewTablet(path)
	defer tablet.Close()
	err = tablet.Open()
	if err != nil {
		t.Errorf("Unable to open tablet: %v", err)
  }
}

