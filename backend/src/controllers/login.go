package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"backend/src/autenticacao"
	"backend/src/banco"
	"backend/src/modelos"
	"backend/src/repositorios"
	"backend/src/respostas"
	"backend/src/seguranca"
)

func Login(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusUnprocessableEntity, "Erro ao ler o conteúdo da chamada: ", erro.Error())
		return
	}

	var usuario modelos.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o conteúdo da chamada: ", erro.Error())
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco de dados: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarioSalvoNoBanco, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao consultar o usuario: ", erro.Error())
		return
	}

	erro = seguranca.VerificarSenha(usuarioSalvoNoBanco.Senha, usuario.Senha)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusUnauthorized, "Usuário não autorizado: ", erro.Error())
		return
	}

	token, erro := autenticacao.CriarToken(usuarioSalvoNoBanco.ID)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao gerar o token: ", erro.Error())
	}

	_, _ = w.Write([]byte(token))
}
