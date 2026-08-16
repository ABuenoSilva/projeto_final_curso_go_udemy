package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/src/config"
	"backend/src/router"
)

func main() {
	config.Carregar()
	r := router.Gerar()
	fmt.Printf("Executando Backend na porta %d", config.Porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
