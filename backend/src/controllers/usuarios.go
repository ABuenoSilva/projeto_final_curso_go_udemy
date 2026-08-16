package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"backend/src/banco"
	"backend/src/modelos"
	"backend/src/repositorios"
	"backend/src/respostas"
)

func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao ler o body: ", erro.Error())
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequest, &usuario); erro != nil {
		respostas.ResponderErro(w, http.StatusUnprocessableEntity, "Erro ao ler o json do body: ", erro.Error())
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario.ID, erro = repositorio.Criar(usuario)

	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao criar o usuário: ", erro.Error())
		return
	}
	respostas.ResponderJSON(w, http.StatusOK, usuario)
}

func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	respostas.ResponderJSON(w, http.StatusOK, "Buscando Usuários!")
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	respostas.ResponderJSON(w, http.StatusOK, "Buscando Usuário!")
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	respostas.ResponderJSON(w, http.StatusOK, "Atualizando Usuário!")
}

func ExcluirUsuario(w http.ResponseWriter, r *http.Request) {
	respostas.ResponderJSON(w, http.StatusOK, "Excluindo Usuário!")
}
