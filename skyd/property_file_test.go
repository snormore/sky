package skyd

import (
	"bytes"
	"testing"
)

// Encode a property file.
func TestEncode(t *testing.T) {
	p := NewPropertyFile("")
	p.CreateProperty("name", "object", "string")
	p.CreateProperty("salary", "object", "float")
	p.CreateProperty("purchaseAmount", "action", "integer")
	p.CreateProperty("isMember", "action", "boolean")

	// Encode
	buffer := new(bytes.Buffer)
	err := p.Encode(buffer)
	if err != nil {
		t.Fatalf("Unable to encode property file: %v", err)
	}
	expected := "[{\"id\":-2,\"name\":\"isMember\",\"type\":\"action\",\"dataType\":\"boolean\"},{\"id\":-1,\"name\":\"purchaseAmount\",\"type\":\"action\",\"dataType\":\"integer\"},{\"id\":1,\"name\":\"name\",\"type\":\"object\",\"dataType\":\"string\"},{\"id\":2,\"name\":\"salary\",\"type\":\"object\",\"dataType\":\"float\"}]\n"
	if buffer.String() != expected {
		t.Fatalf("Invalid property file encoding:\n%v", buffer.String())
	}
}

// Decode a property file.
func TestDecode(t *testing.T) {
	p := NewPropertyFile("")

	// Decode
	buffer := bytes.NewBufferString("[{\"id\":-1,\"name\":\"purchaseAmount\",\"type\":\"action\",\"dataType\":\"integer\"},{\"id\":1,\"name\":\"name\",\"type\":\"object\",\"dataType\":\"string\"},{\"id\":2,\"name\":\"salary\",\"type\":\"object\",\"dataType\":\"float\"}, {\"id\":-2,\"name\":\"isMember\",\"type\":\"action\",\"dataType\":\"boolean\"}]\n")
	err := p.Decode(buffer)
	if err != nil {
		t.Fatalf("Unable to decode property file: %v", err)
	}
	assertProperty(t, p.properties[-2], -2, "isMember", "action", "boolean")
	assertProperty(t, p.properties[-1], -1, "purchaseAmount", "action", "integer")
	assertProperty(t, p.properties[1], 1, "name", "object", "string")
	assertProperty(t, p.properties[2], 2, "salary", "object", "float")
}

