package respostas

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ResponderJSON envia qualquer dado formatado como JSON com o status HTTP desejado
func ResponderJSON(w http.ResponseWriter, status int, dados any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if dados != nil {
		texto, ok := dados.(string)
		if ok {
			dados = struct {
				Mensagem string `json:"mensagem"`
			}{Mensagem: texto}
		}
		if erro := json.NewEncoder(w).Encode(dados); erro != nil {
			ResponderErro(w, http.StatusInternalServerError, "Erro ao gerar o json do usuário:", erro.Error())
		}
	}
}

// ResponderErro envia uma mensagem de erro padronizada em formato JSON
func ResponderErro(w http.ResponseWriter, status int, mensagem string, erro string) {
	mensagem = fmt.Sprintf(mensagem+"%s", erro)
	ResponderJSON(w, status, map[string]string{
		"erro": mensagem,
	})
}
