package torrentFile

import (
	"Gorrent/bencode"
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type torrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func printInfoMap(i map[string]interface{}) {
	fmt.Println("//////// INFO MAP /////////")
	for k, v := range i {
		fmt.Printf("K: %s, V: %s\n", k, v)
	}
	fmt.Println("//////// INFO MAP /////////")
}

func checkAllFieldsExist(data map[string]interface{}) error {
	_, announceExists := data["announce"]
	info, infoExists := data["info"]
	if !announceExists {
		return errors.New("announce field does not exist in top level dictionary")
	} else if !infoExists {
		return errors.New("info field does not exist in top level dictionary")
	}
	infoMap, ok := info.(map[string]interface{})
	printInfoMap(infoMap)
	if !ok {
		return errors.New("info field is not of type map[string]interface{}")
	}
	var errString strings.Builder
	errString.WriteString("The following fields are missing from the metaFileInfo:\n")
	_, lengthExists := infoMap["length"]
	_, nameExists := infoMap["name"]
	_, pieceLengthExists := infoMap["piece length"]
	_, piecesExists := infoMap["pieces"]
	if lengthExists && nameExists && pieceLengthExists && piecesExists {
		return nil
	}
	if !lengthExists {
		errString.WriteString("length\n")
	}
	if !nameExists {
		errString.WriteString("name\n")
	}
	if !pieceLengthExists{
		errString.WriteString("piece length\n")
	}
	if !piecesExists {
		errString.WriteString("pieces\n")
	}
	return errors.New(errString.String())
}

func pieceHashesFormatter(length, pieceLength int, buf []byte) ([][20]byte, error) {
	// If the length fields do not make sense, return an error
	if len(buf) != length || length % pieceLength != 0 || length % 20 != 0 || pieceLength % 20 != 0 {
		return [][20]byte{}, errors.New("got malformed data from torrentFile")
	}
	formattedPieces := make([][20]byte, 0)
	for i := 0; i < length; i += 20 {
		copy(formattedPieces[i][:], buf[i:i+20])
	}
	return formattedPieces, nil
}


func (t *torrentFile) Download(path string) error {
	return nil
}

func NewTorrentFile(path string) (*torrentFile, error) {
	var tf *torrentFile

	// Open torrent file
	file, err := os.Open(path)
	if err != nil {
		return tf, err
	}

	// Defer closing the file safely
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Decode the torrent file using bencoding standard
	decodedData, err := bencode.Decode(bufio.NewReader(io.Reader(file)))
	if err != nil {
		return tf, err
	}

	// Verify that all necessary fields for a torrentFile exist
	if err = checkAllFieldsExist(decodedData); err != nil {
		return tf, err
	}

	// Extract all necessary fields to make a torrentFile struct
	info, _ := decodedData["info"].(map[string]interface{})
	announce := reflect.ValueOf(decodedData["announce"]).String()
	infoHash := sha1.Sum(reflect.ValueOf(decodedData["info"]).Bytes())
	pieceLength := int(reflect.ValueOf(info["piece length"]).Int())
	length := int(reflect.ValueOf(info["length"]).Int())
	name := reflect.ValueOf(info["name"]).String()
	pieces := reflect.ValueOf(info["pieces"]).Bytes()
	formattedPieces, err := pieceHashesFormatter(length, pieceLength, pieces)
	if err != nil {
		return tf, err
	}

	tf = &torrentFile{
		Announce:    announce,
		InfoHash:    infoHash,
		PieceHashes: formattedPieces,
		PieceLength: pieceLength,
		Length:      length,
		Name:        name,
	}
	return tf, nil
}