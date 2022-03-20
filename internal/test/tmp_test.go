package test

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestMatchFiles(t *testing.T) {
	strings, err := filepath.Glob("/data-collection-template/*.template")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(strings)
}
