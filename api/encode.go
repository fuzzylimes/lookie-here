package api

import (
	"bytes"
	"compress/zlib"
	b64 "encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type encode struct {
	URL     string         `json:"url"`
	Markers []encodeMarker `json:"markers"`
}

type encodeMarker struct {
	ID      int    `json:"id"`
	Latlong []int  `json:"latlong"`
	Note    string `json:"note"`
}

func EncodeToBytes(p interface{}) ([]byte, error) {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		return nil, err
	}
	log.Println("uncompressed size (bytes): ", len(buf.Bytes()))
	return buf.Bytes(), nil
}

func Compress(s []byte) []byte {
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write(s); err != nil {
		log.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func B64Encode(s []byte) string {
	return b64.URLEncoding.EncodeToString(s)
}

func EncodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request: ", r.Method)
	if r.Method == "OPTIONS" {
		return
	}
	var p encode
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := EncodeToBytes(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	comp := Compress(b)
	bs4 := B64Encode(comp)
	w.Write([]byte(fmt.Sprintf("{\"d\":\"%s\"}", bs4)))
}
