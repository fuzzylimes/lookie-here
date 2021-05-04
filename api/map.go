package api

import (
	"bytes"
	"compress/zlib"
	b64 "encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Payload struct {
	URL     string   `json:"url"`
	Markers []Marker `json:"markers"`
}

type Marker struct {
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

func DecodeToPayload(s []byte) Payload {
	p := Payload{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
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

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("m") == "encode" {
		log.Println("Got request: ", r.Method)
		if r.Method == "OPTIONS" {
			return
		}
		var p Payload
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
	} else if r.URL.Query().Get("m") == "decode" {
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
	} else {
		http.Error(w, "Invalid API", http.StatusBadRequest)
	}
}
