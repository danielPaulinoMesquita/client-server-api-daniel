package api

import (
	"Client-Server-API/server/database"
	"Client-Server-API/server/service"
	"encoding/json"
	"net/http"
)

type RespBid struct {
	Bid string `json:"bid"`
}

func GetCotacao(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		println("aqui no erro")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cambioDollar := service.RequestCambioDollar()
	model := database.ConvertToModel(cambioDollar)
	database.PersistCambioDollar(model)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := RespBid{Bid: model.Bid}
	json.NewEncoder(w).Encode(resp)
}
