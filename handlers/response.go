package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func sendData(rw http.ResponseWriter, data interface{}, status int) {
	//modifica el estatus en el header
	rw.Header().Set("Contet-Type", "aplication/json")
	rw.WriteHeader(status)

	//Responde al Cliente
	//La data que recibo la convertimos en json
	output, _ := json.Marshal(&data)
	// y el json lo convierte en string
	fmt.Fprintln(rw, string(output))
}

func sendError(rw http.ResponseWriter, status int) {
	//Responde con un status
	rw.WriteHeader(status)
	//Y devuelve el mensaje
	fmt.Fprintln(rw, "Resouece Not Found")
}
