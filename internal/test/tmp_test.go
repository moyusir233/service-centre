package test

import (
	"bytes"
	"encoding/hex"
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

func TestHexEncodeAndDecode(t *testing.T) {
	data := "hello world"

	hexStr1 := fmt.Sprintf("%x", data)
	hexStr2 := hex.EncodeToString([]byte(data))
	if hexStr1 != hexStr2 {
		t.Fatal("hex encode error")
	}

	decodeData1, err := hex.DecodeString(hexStr1)
	if err != nil {
		t.Fatal(err)
	}
	var decodeData2 []byte
	_, err = fmt.Sscanf(hexStr2, "%x", &decodeData2)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(decodeData1, decodeData2) {
		t.Fatal("hex decode error")
	} else if !bytes.Equal([]byte(data), decodeData1) {
		t.Fatal("hex decode error")
	}
}
