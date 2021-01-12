package bencode

import (
	"fmt"
	"strings"
	"testing"
)

func compareEncodings(data map[string]interface{}, result, reference []byte) {
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
	for i, v := range result {
		if v != reference[i] {
			errString.WriteString(fmt.Sprintf("Wrong encoding at index %d\n", i))
			errString.WriteString(fmt.Sprintf("Expected %c, got %c\n", v, reference[i]))
			fmt.Println(errString.String())
			panic("encoding is incorrect")
		}
	}
}

func TestEncode1(t *testing.T) {
	data := make(map[string]interface{})
	data["spam"] = []interface{}{"a", "b"}
	encodedData, err := Encode(data)
	if encodedData == nil || err != nil {
		panic(fmt.Sprintf("Failed to encode data: \n%e", err))
	}
	properEncoding := []byte("d4:spaml1:a1:bee")
	compareEncodings(data, encodedData, properEncoding)
}

func TestEncode2(t *testing.T) {
	data := make(map[string]interface{})
	data["cow"] = "moo"
	data["spam"] = "eggs"
	encodedData, err := Encode(data)
	if encodedData == nil || err != nil {
		panic(fmt.Sprintf("Failed to encode data: \n%e", err))
	}
	properEncoding := []byte("d3:cow3:moo4:spam4:eggse")
	compareEncodings(data, encodedData, properEncoding)
}

func TestEncode3(t *testing.T) {
	data := make(map[string]interface{})
	data["publisher"] = "bob"
	data["publisher-webpage"] = "www.example.com"
	data["publisher.location"] = "home"
	encodedData, err := Encode(data)
	if encodedData == nil || err != nil {
		panic(fmt.Sprintf("Failed to encode data: \n%e", err))
	}
	properEncoding := []byte("d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee")
	compareEncodings(data, encodedData, properEncoding)
}

func TestEncode4(t *testing.T) {
	data := make(map[string]interface{})
	encodedData, err := Encode(data)
	if encodedData == nil || err != nil {
		panic(fmt.Sprintf("Failed to encode data: \n%e", err))
	}
	properEncoding := []byte("de")
	compareEncodings(data, encodedData, properEncoding)
}