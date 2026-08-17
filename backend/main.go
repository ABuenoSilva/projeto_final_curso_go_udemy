package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/src/config"
	"backend/src/router"
)

/* Geração de chave aleatória de 64 bits usada no inicio
func init() {
	chave := make([]byte, 64)
	if _, erro := rand.Read(chave); erro != nil {
		log.Fatal(erro)
	}
	fmt.Println(chave)

	stringBase64 := base64.StdEncoding.EncodeToString(chave)
	fmt.Println(stringBase64)
}
*/

func main() {
	config.Carregar()
	r := router.Gerar()
	fmt.Printf("Executando Backend na porta %d", config.Porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Porta), r))
}
