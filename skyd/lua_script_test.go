package skyd

import (
  "testing"
)

// Ensure that the lua script can initialize.
func TestLuaScriptInit(t *testing.T) {
  p := NewPropertyFile("")
  l := NewLuaScript(p, "function test() print('hello') end\n")
  l.Execute("test", 0)
  l.Destroy()
}

// Ensure that the lua script can extract event property references.
func TestLuaScriptExtractPropertyReferences(t *testing.T) {
  p := NewPropertyFile("")
  p.CreateProperty("name", "object", "string")
  p.CreateProperty("salary", "object", "float")
  p.CreateProperty("purchaseAmount", "action", "integer")
  p.CreateProperty("isMember", "action", "boolean")

  l := NewLuaScript(p, "x = event:name() if event.salary > 100 then print(event.purchaseAmount, event, event:name()) end")
  properties, err := l.ExtractPropertyReferences()
  if err != nil {
    t.Fatalf("Unable to extract properties: %v", err)
  }
  if len(properties) != 3 {
    t.Fatalf("Expected %v properties, got %v", 3, len(properties))
  }
  if p.GetPropertyByName("purchaseAmount") != properties[0] {
    t.Fatalf("Expected %v, got %v", p.GetPropertyByName("purchaseAmount"), properties[0])
  }
  if p.GetPropertyByName("name") != properties[1] {
    t.Fatalf("Expected %v, got %v", p.GetPropertyByName("name"), properties[1])
  }
  if p.GetPropertyByName("salary") != properties[2] {
    t.Fatalf("Expected %v, got %v", p.GetPropertyByName("salary"), properties[2])
  }
}

// Ensure that the lua script can generate a header.
/*
func TestLuaScriptGenerateHeader(t *testing.T) {
  p := NewPropertyFile("")
  p.CreateProperty("name", "object", "string")
  p.CreateProperty("salary", "object", "float")
  p.CreateProperty("purchaseAmount", "action", "integer")
  p.CreateProperty("isMember", "action", "boolean")

  l := NewLuaScript(p, "function test(event) x = event:name() if event.salary > 100 then print(event.purchaseAmount, event, event:name()) end end")
  err := l.Init()
  if err != nil {
    t.Fatalf("Unable to initialize: %v", err)
  }

  fmt.Println("--START HEADER--")
  fmt.Println(l.Header())
  fmt.Println("--END HEADER--")
}

*/