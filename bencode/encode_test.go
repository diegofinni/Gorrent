package bencode

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestNewEncoder(t *testing.T) {
	e := NewEncoder()
	if e == nil {
		panic(e)
	}
}

func TestEncoder_Encode(t *testing.T) {
	encoder := NewEncoder()
	data := make(map[string]interface{})
	data["cow"] = "moo"
	data["spam"] = []interface{}{"a", "b"}
	encoder.Encode(data)
	encodedData, err := encoder.ReadAll()
	if encodedData == nil || err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(encodedData)
	decodedData, err := Decode(bufio.NewReader(buf))
	if decodedData == nil || err != nil {
		panic(err)
	}
	originalLen, _ := fmt.Println(data)
	decodedLen,  _ := fmt.Println(decodedData)
	if originalLen != decodedLen {
		panic("decoded data does not equal original data")
	}
}

func TestEncoder_ReadAll(t *testing.T) {

}