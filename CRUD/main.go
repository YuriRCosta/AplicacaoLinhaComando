package main

import (
	"crud/servidor"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const URL_USUARIO_ID = "/usuarios/{id}"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/usuarios", servidor.CriarUsuario).Methods(http.MethodPost)
	router.HandleFunc("/usuarios", servidor.BuscarUsuarios).Methods(http.MethodGet)
	router.HandleFunc(URL_USUARIO_ID, servidor.BuscarUsuario).Methods(http.MethodGet)
	router.HandleFunc(URL_USUARIO_ID, servidor.AtualizarUsuario).Methods(http.MethodPut)
	router.HandleFunc(URL_USUARIO_ID, servidor.DeletarUsuario).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":5000", router))
}
