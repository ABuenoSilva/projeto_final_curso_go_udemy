package middlewares

import (
	"fmt"
	"log"
	"net/http"

	"backend/src/autenticacao"
	"backend/src/respostas"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(" ")
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

func Autenticar(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if erro := autenticacao.ValidarToken(r); erro != nil {
			respostas.ResponderErro(w, http.StatusUnauthorized, "Token inválido: ", erro.Error())
			return
		}
		next(w, r)
	}
}
