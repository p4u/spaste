package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

func update(d *leveldb.DB, k, v string) error {
	err := d.Put([]byte(k), []byte(v), nil)
	return err
}

func get(d *leveldb.DB, k string) []byte {
	data, err := d.Get([]byte(k), nil)
	if err != nil {
		return nil
	}
	return data
}

func delete(d *leveldb.DB, k string) error {
	err := d.Delete([]byte(k), nil)
	return err
}

func dumpDB(d *leveldb.DB, w http.ResponseWriter) {
	iter := d.NewIterator(nil, nil)
	for iter.Next() {
		k := iter.Key()
		w.Write([]byte(fmt.Sprintf("%s\n", k)))
	}
	iter.Release()
}

func reply(resp string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(resp))
}

func handler(w http.ResponseWriter, req *http.Request) {
	var resp string
	params := strings.Split(req.RequestURI, "/")
	log.Printf("Received %s from %s\n", params, req.RemoteAddr)
	if params[1] == "add" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		data := buf.String()
		var key string
		if len(params) > 2 {
			key = params[2]
		}
		if len(key) < 1 {
			key = fmt.Sprintf("%d", time.Now().Unix())
		}
		log.Printf("Adding %d bytes of data to %s\n", len(data), key)
		update(db, key, data)
		resp = "key=" + key
	}
	if params[1] == "get" && len(params) > 2 {
		log.Printf("Get %s\n", params[2])
		resp = string(get(db, params[2]))
	}
	if params[1] == "list" {
		dumpDB(db, w)
	}
	reply(resp, w)
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Usage: " + os.Args[0] + " <port>")
		os.Exit(2)
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	db, err = leveldb.OpenFile("db.spaste", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 4 * time.Second,
		ReadTimeout:       4 * time.Second,
		WriteTimeout:      4 * time.Second,
		IdleTimeout:       3 * time.Second,
	}

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
	http.HandleFunc("/add/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})

	log.Printf("Starting server ar port %d\n", port)
	srv.SetKeepAlivesEnabled(false)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

}
