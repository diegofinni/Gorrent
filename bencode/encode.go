package bencode

import (
	"bytes"
	"errors"
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

func (e *encoder) encodeUint(i uint64) {
	e.buf.WriteByte('i')
	e.buf.WriteString(strconv.FormatUint(i, 10))
	e.buf.WriteByte('e')
}

func (e *encoder) encodeString(str string) {
	e.buf.WriteString(strconv.FormatInt(int64(len(str)), 10))
	e.buf.WriteByte(':')
	e.buf.WriteString(str)
}

func (e *encoder) encodeList(lst []interface{}) error {
	var err error
	e.buf.WriteByte('l')
	for _, v := range lst {
		err = e.encodeByType(v)
		if err != nil {
			return err
		}
	}
	e.buf.WriteByte('e')
	return nil
}

func (e *encoder) encodeDictionary(dic map[string]interface{}) error {
	var err error
	e.buf.WriteByte('d')
	for k, v := range dic {
		err = e.encodeByType(k)
		if err != nil {
			return err
		}
		err = e.encodeByType(v)
		if err != nil {
			return err
		}
	}
	e.buf.WriteByte('e')
	return nil
}


// Bencode the given item and add it to the encoder buffer
func (e *encoder) encodeByType(item interface{}) error {
	var err error
	switch v := item.(type) {
	case string:
		e.encodeString(reflect.ValueOf(v).String())
	case int, int8, int16, int32, int64:
		e.encodeInt(reflect.ValueOf(v).Int())
	case uint, uint8, uint16, uint32, uint64:
		e.encodeUint(reflect.ValueOf(v).Uint())
	case map[string]interface{}:
		err = e.encodeDictionary(v)
	case []interface{}:
		err = e.encodeList(v)
	default:
		err = errors.New("item type did not match any of the supported ones")
	}
	return err
}

// Return a byte array of all encoded objects given to the encoder
func (e *encoder) readAll() ([]byte, error) {
	returnBuf := make([]byte, e.buf.Len())
	_, err := e.buf.Read(returnBuf)
	if err != nil {
		return nil, err
	}
	return returnBuf, nil
}

func Encode(data interface{}) ([]byte, error) {
	e := &encoder{*bytes.NewBuffer(make([]byte, 0))}
	dataType := reflect.TypeOf(data).Kind()
	if dataType != reflect.Map {
		return nil, errors.New("data type must be Map")
	}
	err := e.encodeByType(data)
	if err != nil {
		return nil, err
	}
	return e.readAll()
}
