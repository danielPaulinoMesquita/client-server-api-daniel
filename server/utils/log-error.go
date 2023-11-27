package utils

import "log"

const LOG_ERROR_DB = "Error happens when persists in the database"
const LOG_ERROR_REQUEST_CAMBIO = "Error to request api economia"

func LogError(err error, message string) {
	log.Println(message, err)
}

func TreatError(err error) {
	log.Fatalln(err)
}
