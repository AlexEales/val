package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	ErrorNoSuchKey = errors.New("no such key")

	store = make(map[string]string)
)

func Delete(key string) error {
	delete(store, key)

	return nil
}

func Get(key string) (string, error) {
	value, ok := store[key]
	if !ok {
		return "", fmt.Errorf("%w: %q", ErrorNoSuchKey, key)
	}

	return value, nil
}

func Put(key, value string) error {
	store[key] = value
	return nil
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Get(key)
	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(value))
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", putHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", getHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", deleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
