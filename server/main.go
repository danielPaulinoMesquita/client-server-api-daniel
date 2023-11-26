package main

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"io"
	"log"
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

type CambioDollarModel struct {
	Id         int `gorm:"primaryKey"`
	Code       string
	Codein     string
	Name       string
	High       string
	Low        string
	VarBid     string
	PctChange  string
	Bid        string
	Ask        string
	Timestamp  string
	CreateDate string
	gorm.Model
}

type RespBid struct {
	Bid string `json:"bid"`
}

// todo study later
// sudo docker-compose -f docker-compose-sqlite.yaml up -d
// https://medium.com/@jamal.kaksouri/the-complete-guide-to-context-in-golang-efficient-concurrency-management-43d722f6eaea
func main() {
	http.HandleFunc("/cotacao", getCotacao)
	println("teste")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
		return
	}
}

func convertToModel(body *RespDollarBody) CambioDollarModel {
	return CambioDollarModel{
		Code:       body.USDBRL.Code,
		Codein:     body.USDBRL.Codein,
		Name:       body.USDBRL.Name,
		High:       body.USDBRL.High,
		Low:        body.USDBRL.Low,
		VarBid:     body.USDBRL.VarBid,
		PctChange:  body.USDBRL.PctChange,
		Bid:        body.USDBRL.Bid,
		Ask:        body.USDBRL.Ask,
		Timestamp:  body.USDBRL.Timestamp,
		CreateDate: body.USDBRL.CreateDate,
	}
}

func requestCambioDollar() *RespDollarBody {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil)

	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// FIRST LOG
		log.Println("Error to request api economia", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var c RespDollarBody
	err = json.Unmarshal(body, &c)

	if err != nil {
		panic(err)
	}

	return &c
}

func persistCambioDollar(model CambioDollarModel) {
	contextDatabase, cancel := context.WithTimeout(context.Background(), time.Millisecond*400) // <- change time to 10ms
	defer cancel()

	//db, err := gorm.Open(
	//	mysql.Open("root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"),
	//	&gorm.Config{})
	db, err := sql.Open("sqlite3", "clientServerApi.db")

	if err != nil {
		panic(err)
	}
	//
	//err = db.(&CambioDollarModel{})
	//if err != nil {
	//	panic(err)
	//}

	createTableSQL := `CREATE TABLE IF NOT EXISTS cambio_dollar_models (
		id INTEGER PRIMARY KEY,
		code TEXT,
		codein TEXT,
		name TEXT, 
		high TEXT, 
		low TEXT, 
		var_bid TEXT, 
		pct_change TEXT, 
		bid TEXT, ask TEXT, 
		timestamp TEXT, create_date TEXT,
		 created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
		 updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}

	//err = db.PrepareContext(contextDatabase).Create(&model).Error
	stmt, err := db.PrepareContext(contextDatabase, `INSERT INTO 
    	cambio_dollar_models (code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	if err != nil {
		panic(err)
	}

	_, err = stmt.ExecContext(contextDatabase,
		model.Code,
		model.Codein,
		model.Name,
		model.High,
		model.Low,
		model.VarBid,
		model.PctChange,
		model.Bid,
		model.Ask,
		model.Timestamp,
		model.CreateDate)

	if err != nil {
		// SECOND LOG
		log.Println("Error happens when persists in the database", err)
	}

}

func getCotacao(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		println("aqui no erro")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	println("here -----------")

	cambioDollar := requestCambioDollar()
	model := convertToModel(cambioDollar)
	persistCambioDollar(model)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := RespBid{Bid: model.Bid}
	json.NewEncoder(w).Encode(resp)
}
