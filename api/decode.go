package api

import (
	"bytes"
	"compress/zlib"
	b64 "encoding/base64"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type decode struct {
	URL     string         `json:"url"`
	Markers []decodeMarker `json:"markers"`
}

type decodeMarker struct {
	ID      int    `json:"id"`
	Latlong []int  `json:"latlong"`
	Note    string `json:"note"`
}

func B64Decode(s []byte) []byte {
	res, err := b64.URLEncoding.DecodeString(string(s))
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func Decompress(s []byte) []byte {
	zr, err := zlib.NewReader(bytes.NewReader(s))
	if err != nil {
		log.Fatal(err)
	}
	result, err := ioutil.ReadAll(zr)
	if err != nil {
		log.Fatal(err)
	}
	if err := zr.Close(); err != nil {
		log.Fatal(err)
	}
	return result
}

func DecodeToPayload(s []byte) decode {
	p := decode{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func DecodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("New Request")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Got: ", string(b))
	enc := B64Decode(b)
	decomp := Decompress(enc)
	p := DecodeToPayload(decomp)
	jsonData, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
