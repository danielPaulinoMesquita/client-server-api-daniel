package service

import (
	"Client-Server-API/server/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type RespDollarBody struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	}
}

func RequestCambioDollar() *RespDollarBody {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil)

	if err != nil {
		utils.TreatError(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// FIRST LOG
		utils.LogError(err, utils.LOG_ERROR_REQUEST_CAMBIO)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.TreatError(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.TreatError(err)
	}

	var c RespDollarBody
	err = json.Unmarshal(body, &c)

	if err != nil {
		utils.TreatError(err)
	}

	return &c
}
