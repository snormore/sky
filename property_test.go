package skyd

import (
	"testing"
)

// Ensure that property names can use valid characters.
func TestPropertyName(t *testing.T) {
	if _, err := NewProperty(0, "  Property_no - 2", false, "string"); err != nil {
		t.Fatal("Property name:", err)
	}
}

// Ensure that property names cannot have illegal characters.
func TestPropertyNameCannotContainInvalidCharacters(t *testing.T) {
	if _, err := NewProperty(0, "yes\\no", false, "string"); err.Error() != "Property name contains invalid characters: yes\\no" {
		t.Fatal("Invalid name:", err)
	}
}

