package bencode

import (
	"bytes"
	"reflect"
	"strconv"
)
type encoder struct {
	buf bytes.Buffer
}

func (e *encoder) encodeInt(i int64) {
	e.buf.WriteByte('i')
	e.buf.WriteString(strconv.FormatInt(i, 10))
	e.buf.WriteByte('e')
}

func (e *encoder) encodeString(str string) {
	e.buf.WriteString(strconv.FormatInt(int64(len(str)), 10))
	e.buf.WriteByte(':')
	e.buf.WriteString(str)
}

func (e *encoder) encodeList(lst []interface{}) {
	e.buf.WriteByte('l')
	for _, v := range lst {
		e.Encode(v)
	}
	e.buf.WriteByte('e')
}

func (e *encoder) encodeDictionary(dic map[string]interface{}) {
	e.buf.WriteByte('d')
	for k, v := range dic {
		e.Encode(k)
		e.Encode(v)
	}
	e.buf.WriteByte('e')
}

// Return a new encoder with an empty buffer
func NewEncoder() *encoder {
	return &encoder{
		*bytes.NewBuffer(make([]byte, 0)),
	}
}

// Bencode the given item and add it to the encoder buffer
func (e *encoder) Encode(item interface{}) {
	tp := reflect.TypeOf(item)
	switch tp.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.encodeInt(reflect.ValueOf(item).Int())
	case reflect.String:
		e.encodeString(reflect.ValueOf(item).String())
	case reflect.Map:
		e.encodeDictionary(item.(map[string]interface{}))
	case reflect.Array, reflect.Slice:
		e.encodeList(item.([]interface{}))
	}
}

// Return a byte array of all encoded objects given to the encoder
func (e *encoder) ReadAll() ([]byte, error) {
	returnBuf := make([]byte, e.buf.Len())
	_, err := e.buf.Read(returnBuf)
	if err != nil {
		return nil, err
	}
	return returnBuf, nil
}