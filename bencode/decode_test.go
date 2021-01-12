package bencode

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func compareDecodings(data []byte, result, reference map[string]interface{}) {
	var errString strings.Builder
	errString.WriteString("Data: ")
	errString.WriteString(fmt.Sprintln(data))
	errString.WriteString(fmt.Sprintf("Result: %s\n", result))
	errString.WriteString(fmt.Sprintf("Reference: %s\n", reference))
	if len(result) != len(reference) {
		errString.WriteString(fmt.Sprintf("Result len: %d, Reference len: %d", len(result), len(reference)))
		fmt.Println(errString.String())
		panic("encoding is wrong length")
	}
	for k := range result {
		_, exists := reference[k]
		if !exists {
			errString.WriteString(fmt.Sprintf("Key'%s' does not exist in reference decoding", k))
			fmt.Println(errString.String())
			panic("decoding is incorrect")
		}
	}
}

func TestDecode1(t *testing.T) {
	data := []byte("d4:spaml1:a1:bee")
	decodedData, err := Decode(bufio.NewReader(bytes.NewBuffer(data)))
	if err != nil {
		panic(err)
	}
	properDecoding := make(map[string]interface{})
	properDecoding["spam"] = []interface{}{"a", "b"}
	compareDecodings(data, decodedData, properDecoding)
}

func TestDecode2(t *testing.T) {
	data := []byte("d3:cow3:moo4:spam4:eggse")
	decodedData, err := Decode(bufio.NewReader(bytes.NewBuffer(data)))
	if err != nil {
		panic(err)
	}
	properDecoding := make(map[string]interface{})
	properDecoding["cow"] = "moo"
	properDecoding["spam"] = "eggs"
	compareDecodings(data, decodedData, properDecoding)
}

func TestDecode3(t *testing.T) {
	data := []byte("d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee")
	decodedData, err := Decode(bufio.NewReader(bytes.NewBuffer(data)))
	if err != nil {
		panic(err)
	}
	properDecoding := make(map[string]interface{})
	properDecoding["publisher"] = "bob"
	properDecoding["publisher-webpage"] = "www.example.com"
	properDecoding["publisher.location"] = "home"
	compareDecodings(data, decodedData, properDecoding)
}

func TestDecode4(t *testing.T) {
	data := []byte("de")
	decodedData, err := Decode(bufio.NewReader(bytes.NewBuffer(data)))
	if err != nil {
		panic(err)
	}
	properDecoding := make(map[string]interface{})
	compareDecodings(data, decodedData, properDecoding)
}