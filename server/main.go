package main

import (
	"Client-Server-API/server/api"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func main() {
	http.HandleFunc("/cotacao", api.GetCotacao)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
		return
	}
}
