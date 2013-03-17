package skyd

import (
	"testing"
)

// Ensure that the lua script can initialize.
func TestExecutionEngineInit(t *testing.T) {
	p := NewPropertyFile("")
	l, _ := NewExecutionEngine(p, "function test() print('hello') end\n")
	l.Destroy()
}

// Ensure that the lua script can extract event property references.
func TestExecutionEngineExtractPropertyReferences(t *testing.T) {
	p := NewPropertyFile("")
	p.CreateProperty("name", "object", "string")
	p.CreateProperty("salary", "object", "float")
	p.CreateProperty("purchaseAmount", "action", "integer")
	p.CreateProperty("isMember", "action", "boolean")

	l, err := NewExecutionEngine(p, "function f(event) x = event:name() if event.salary > 100 then print(event.purchaseAmount, event, event:name()) end end")
	if err != nil {
		t.Fatalf("Unable to create execution engine: %v", err)
	}
	if len(l.propertyRefs) != 3 {
		t.Fatalf("Expected %v properties, got %v", 3, len(l.propertyRefs))
	}
	if p.GetPropertyByName("purchaseAmount") != l.propertyRefs[0] {
		t.Fatalf("Expected %v, got %v", p.GetPropertyByName("purchaseAmount"), l.propertyRefs[0])
	}
	if p.GetPropertyByName("name") != l.propertyRefs[1] {
		t.Fatalf("Expected %v, got %v", p.GetPropertyByName("name"), l.propertyRefs[1])
	}
	if p.GetPropertyByName("salary") != l.propertyRefs[2] {
		t.Fatalf("Expected %v, got %v", p.GetPropertyByName("salary"), l.propertyRefs[2])
	}
}

// Ensure that the lua script can generate a header.
/*
func TestExecutionEngineGenerateHeader(t *testing.T) {
  p := NewPropertyFile("")
  p.CreateProperty("name", "object", "string")
  p.CreateProperty("salary", "object", "float")
  p.CreateProperty("purchaseAmount", "action", "integer")
  p.CreateProperty("isMember", "action", "boolean")

  l := NewExecutionEngine(p, "function test(event) x = event:name() if event.salary > 100 then print(event.purchaseAmount, event, event:name()) end end")
  err := l.Init()
  if err != nil {
    t.Fatalf("Unable to initialize: %v", err)
  }

  fmt.Println(l.FullAnnotatedSource())
}

*/
