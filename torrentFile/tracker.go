package torrentFile

import (
	"Gorrent/peers"
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

const port = 6882

func makePeerID() [20]byte {
	unixNanoTime := time.Now().UnixNano()
	timeInBytes := reflect.ValueOf(unixNanoTime).Bytes()
	return sha1.Sum(timeInBytes)
}

func (tf *torrentFile) sendTrackerGetRequest() ([]peers.Peer, error) {
	// Create base url from announce link
	baseURL, err := url.Parse(tf.Announce)
	if err != nil {
		return []peers.Peer{}, nil
	}
	// Create peer id and query parameters
	id := makePeerID()
	query := baseURL.Query()
	query.Add("info_hash", string(tf.InfoHash[:]))
	query.Add("peer_id", string(id[:]))
	query.Add("port", strconv.Itoa(port))
	query.Add("uploaded", "0")
	query.Add("downloaded", "0")
	query.Add("left", strconv.Itoa(tf.Length))
	query.Add("compact", "1")
	// Add parameters to query and return
	baseURL.RawQuery = query.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		return []peers.Peer{}, err
	}
	fmt.Println(resp)
	return nil, nil
}