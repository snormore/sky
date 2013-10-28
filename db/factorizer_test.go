package db

import (
	"io/ioutil"
	"os"
	"testing"
)

// Ensure that we can factorize and defactoize values
func TestFactorization(t *testing.T) {
	withFactorizer(func(f *factorizer) {
		num, err := f.Factorize("foo", "bar", "/index.html", true)
		if err != nil || num != 1 {
			t.Fatalf("Wrong factorization: exp: %v, got: %v (%v)", 1, num, err)
		}
		num, err = f.Factorize("foo", "bar", "/about.html", true)
		if err != nil || num != 2 {
			t.Fatalf("Wrong factorization: exp: %v, got: %v (%v)", 2, num, err)
		}

		str, err := f.Defactorize("foo", "bar", 1)
		if err != nil || str != "/index.html" {
			t.Fatalf("Wrong defactorization: exp: %v, got: %v (%v)", "/index.html", str, err)
		}
		str, err = f.Defactorize("foo", "bar", 2)
		if err != nil || str != "/about.html" {
			t.Fatalf("Wrong defactorization: exp: %v, got: %v (%v)", "/about.html", str, err)
		}
	})
}

func withFactorizer(fn func(f *factorizer)) {
	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)

	f := NewFactorizer(path, false, 4096, 126).(*factorizer)
	if err := f.Open(); err != nil {
		panic(err.Error())
	}
	defer f.Close()

	fn(f)
}
