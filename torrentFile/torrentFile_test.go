package torrentFile

import (
	"testing"
)

func TestNewTorrentFile(t *testing.T) {
	path := "/Users/diegosanmiguel/Documents/goWorkSpace/src/Gorrent/testingFiles/debian.torrent"
	_, err := NewTorrentFile(path)
	if err != nil {
		panic(err)
	}
}