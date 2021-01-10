package bencode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"unicode"
)

func decodeInt(reader *bufio.Reader) (int64, error) {
	// Read the very first byte, which should be 'i', then proceed
	_, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	// Read up until you find the character 'e'
	bytes, err := reader.ReadBytes('e')
	if err != nil {
		return 0, err
	}
	// Cast the byte array to a string then parse it for an int64
	return strconv.ParseInt(string(bytes[:len(bytes)-1]), 10, 64)
}

func decodeString(reader *bufio.Reader) (string, error) {
	// Read up until the character ':'
	bytes, err := reader.ReadBytes(':')
	if err != nil {
		return "", err
	}

	// Convert the byte array to string then parse it to get an int64
	length, err := strconv.ParseInt(string(bytes[:len(bytes)-1]), 10, 64)
	if err != nil {
		return "", err
	} else if length < 0 {
		return "", errors.New("string length is negative")
	}

	// Read length number of bytes and return that byte array cast to string
	buf := make([]byte, length)
	bytesRead, err := io.ReadFull(reader, buf)
	if err != nil {
		return "", err
	} else if int64(bytesRead) != length {
		return "", errors.New("determined length of decoded string was incorrect")
	}
	return string(buf), nil
}

func decodeList(reader *bufio.Reader) ([]interface{}, error) {
	_, _ = reader.ReadByte()
	lst := make([]interface{}, 0)
	for {
		// The next byte being 'e' indicates that the dictionary has ended
		nextByte, err := reader.Peek(1)
		if err != nil {
			return lst, err
		} else if nextByte[0] == 'e' {
			return lst, nil
		}

		// Decode next element in the lst
		value, err := decodeNextElement(reader)
		if err != nil {
			return lst, err
		}
		lst = append(lst, value)
	}
}

func decodeDictionary(reader *bufio.Reader) (map[string]interface{}, error) {
	_, _ = reader.ReadByte()
	dic := make(map[string]interface{})
	for {
		// The next byte being 'e' indicates that the dictionary has ended
		nextByte, err := reader.Peek(1)
		if err != nil {
			return dic, err
		} else if nextByte[0] == 'e' {
			return dic, nil
		}

		// Decode the key for the next kv pair
		key, err := decodeString(reader)
		if err != nil {
			return dic, err
		}

		// Decode the value corresponding to the key
		value, err := decodeNextElement(reader)
		if err != nil {
			return dic, err
		}

		// Place decoded key and value into dictionary
		dic[key] = value
	}
}

func decodeNextElement(reader *bufio.Reader) (interface{}, error) {
	// Determine what the type of the next value will be by peeking the first byte
	typeByte, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}

	// Call the proper decode function by casing on the value of typeByte
	var value interface{}
	switch typeByte[0] {
	case 'd':
		value, err = decodeDictionary(reader)
	case 'i':
		value, err = decodeInt(reader)
	case 'l':
		value, err = decodeList(reader)
	default:
		if unicode.IsNumber(rune(typeByte[0])) {
			value, err = decodeString(reader)
		} else {
			return nil, errors.New("data was not encoded to proper bencode standard")
		}
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func Decode(reader *bufio.Reader) (map[string]interface{}, error) {
	firstByte, err := reader.Peek(1)
	if err != nil {
		return nil, err
	}
	if firstByte[0] == 'd' {
		return decodeDictionary(reader)
	} else {
		return nil, errors.New("bencode data must start with a dictionary")
	}
}