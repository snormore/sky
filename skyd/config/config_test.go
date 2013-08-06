package config

import (
	"bytes"
	"testing"
)

const testConfigFileA = `
port=9000
data-path="/home/data"
pid-path = "/home/pid"
nosync = true
`

// Decode a configuration file.
func TestDecode(t *testing.T) {
	config := NewConfig()
	err := config.Decode(bytes.NewBufferString(testConfigFileA))

	if err != nil {
		t.Fatalf("Unable to decode: %v", err)
	} else if config.Port != 9000 {
		t.Fatalf("Invalid port: %v", config.Port)
	} else if config.DataPath != "/home/data" {
		t.Fatalf("Invalid data path: %v", config.DataPath)
	} else if config.PidPath != "/home/pid" {
		t.Fatalf("Invalid pid path: %v", config.PidPath)
	} else if config.NoSync != true {
		t.Fatalf("Invalid nosync option: %v", config.NoSync)
	}
}
