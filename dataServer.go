package main

import (
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"errors"
	"strings"
)

var hostIP string

func upload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	remoteAddr := strings.Split(r.RemoteAddr, ":")
	remoteIP := remoteAddr[0]

	if remoteIP != hostIP {
		fmt.Printf("Illegal request from %s, should be %s", remoteIP, hostIP)
	}

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
		fmt.Printf("%s is wrote.\n", fileName)

		return
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	remoteAddr := strings.Split(r.RemoteAddr, ":")
	remoteIP := remoteAddr[0]

	if remoteIP != hostIP {
		fmt.Printf("Illegal request from %s, should be %s", remoteIP, hostIP)
	}

	if r.Method == "GET" {
		fileName := r.Header.Get("Filename")
		if fileName == "" {
			w.WriteHeader(500)

			/*for a:=0; a<10; a++ {
				w.Write([]byte("hello"))
				if f,ok := w.(http.Flusher); ok {
					f.Flush()
				}

				time.Sleep(1000000000)
				fmt.Println("write")

			}*/
			return
		}
		s, err := ioutil.ReadFile("data/" + fileName)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(s)
		fmt.Printf("%s is read.\n", fileName)

		return
	}
}

func register(addr string) error {
	url := "http://" + addr + ":8081/register"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resText := string(body)

	if resText != "Success" {
		return errors.New(fmt.Sprintf("%s %s\n", url, resText))
	}
	fmt.Printf("%s %s\n", url, resText)

	return nil
}

func main() {
	temp := flag.String("addr", "", "Name server address")
	flag.Parse()
	hostIP = *temp
	if hostIP == "" {
		fmt.Println("Need address.")
		return
	}
	err := register(hostIP)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Register failed, exiting")
		return
	}
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/download", download)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
