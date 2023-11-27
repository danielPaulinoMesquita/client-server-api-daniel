package database

import (
	"Client-Server-API/server/service"
	"Client-Server-API/server/utils"
	"context"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

const createTable = `CREATE TABLE IF NOT EXISTS cambio_dollar_models (
	id INTEGER PRIMARY KEY,
	code TEXT,
	codein TEXT,
	name TEXT, 
	high TEXT, 
	low TEXT, 
	var_bid TEXT, 
	pct_change TEXT, 
	bid TEXT, ask TEXT, 
	timestamp TEXT, 
	create_date TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

const insertTable = `INSERT INTO cambio_dollar_models (
	code, codein, name,high, low, var_bid, pct_change,
    bid,  ask, timestamp,create_date) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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

func ConvertToModel(body *service.RespDollarBody) CambioDollarModel {
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

func PersistCambioDollar(model CambioDollarModel) {
	contextDatabase, cancel := context.WithTimeout(context.Background(), time.Millisecond*400) // <- todo change time to 10ms
	defer cancel()

	db, err := sql.Open("sqlite3", "clientServerApi.db")

	if err != nil {
		utils.TreatError(err)
	}

	_, err = db.Exec(createTable)
	if err != nil {
		utils.TreatError(err)
	}

	stmt, err := db.PrepareContext(contextDatabase, insertTable)

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	if err != nil {
		utils.TreatError(err)
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
		utils.LogError(err, utils.LOG_ERROR_DB)
	}

}
