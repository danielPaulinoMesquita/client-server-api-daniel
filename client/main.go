package main

import (
	"Client-Server-API/client/file"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type RespBid struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao requisitar o valor atual do câmbio!")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var bid RespBid
	err = json.Unmarshal(body, &bid)
	file.WriterFile(bid.Bid)
}
