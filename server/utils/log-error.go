package utils

import "log"

const LOG_ERROR_DB = "Error aconteceu ao tentar persistir no banco de dados"
const LOG_ERROR_REQUEST_CAMBIO = "Error na requisição da api de câmbio"

func LogError(err error, message string) {
	log.Println(message, err)
}

func TreatError(err error) {
	log.Fatalln(err)
}
