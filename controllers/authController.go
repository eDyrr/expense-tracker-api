package controllers

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received a registration request")

	var data map[string]string

}
