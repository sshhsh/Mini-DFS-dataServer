package main

import (
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
)

func upload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "POST" {
		s, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		fileName := r.Header.Get("Filename")
		if fileName == "" {
			w.WriteHeader(500)
			return
		}

		fmt.Println(r.ContentLength)
		err = ioutil.WriteFile("data/"+fileName, s, 0666)
		md5Ctx := md5.New()
		md5Ctx.Write(s)
		md5Result := md5Ctx.Sum(nil)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write([]byte(hex.EncodeToString(md5Result)))

	}
}

func download(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method == "GET" {
		fileName := r.Header.Get("Filename")
		if fileName == "" {
			w.WriteHeader(500)
			return
		}
		s,err := ioutil.ReadFile("data/"+fileName)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(s)
	}
}

func main() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", download)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

