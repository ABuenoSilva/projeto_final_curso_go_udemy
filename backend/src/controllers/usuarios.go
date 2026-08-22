package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"backend/src/autenticacao"
	"backend/src/banco"
	"backend/src/modelos"
	"backend/src/repositorios"
	"backend/src/respostas"
	"backend/src/seguranca"

	"github.com/gorilla/mux"
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

	if erro = usuario.Preparar("cadastro"); erro != nil {
		respostas.ResponderErro(w, http.StatusUnprocessableEntity, "Erro ao preparar o usuário: ", erro.Error())
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
	nomeOuNick := strings.ToLower(r.URL.Query().Get("usuario"))

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar(nomeOuNick)

	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao consultar o banco: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, usuarios)
}

func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o parâmetro: ", erro.Error())
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()
	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)

	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao consultar o banco: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, usuario)
}

func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o parâmetro: ", erro.Error())
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusUnauthorized, "Usuário não autorizado: ", erro.Error())
	}

	if usuarioID != usuarioIDNoToken {
		respostas.ResponderErro(w, http.StatusForbidden, "Operação não permitida: ", "Não é possível atualizar um usuário que não seja o seu")
		return
	}

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

	if erro = usuario.Preparar("edicao"); erro != nil {
		respostas.ResponderErro(w, http.StatusUnprocessableEntity, "Erro ao preparar o usuário: ", erro.Error())
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	erro = repositorio.Atualizar(usuarioID, usuario)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao atualizar o usuário: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, "Usuário atualizado com sucesso!")
}

func ExcluirUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o parâmetro: ", erro.Error())
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusUnauthorized, "Usuário não autorizado: ", erro.Error())
	}

	if usuarioID != usuarioIDNoToken {
		respostas.ResponderErro(w, http.StatusForbidden, "Operação não permitida: ", "Não é possível excluir um usuário que não seja o seu")
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	erro = repositorio.Excluir(usuarioID)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao excluir o usuário: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, "Usuário excluído com sucesso!")
}

func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusUnauthorized, "Usuário não autorizado: ", erro.Error())
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o parâmetro: ", erro.Error())
		return
	}

	if seguidorID == usuarioID {
		respostas.ResponderErro(w, http.StatusForbidden, "Operação não permitida :", "Um usuário não pode seguir ele mesmo")
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Seguir(usuarioID, seguidorID); erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao seguir o usuário: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusCreated, fmt.Sprintf("Seguiu o usuário %d", usuarioID))
}

func PararDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusUnauthorized, "Usuário não autorizado: ", erro.Error())
	}

	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o parâmetro: ", erro.Error())
		return
	}

	if seguidorID == usuarioID {
		respostas.ResponderErro(w, http.StatusForbidden, "Operação não permitida :", "Um usuário não pode deixar de seguir ele mesmo")
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.PararDeSeguir(usuarioID, seguidorID); erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao deixar de seguir o usuário: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, fmt.Sprintf("Deixou de seguir o usuário %d", usuarioID))
}

func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o parâmetro: ", erro.Error())
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	seguidores, erro := repositorio.BuscarSeguidores(usuarioID)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao buscar seguidores: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, seguidores)
}

func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o parâmetro: ", erro.Error())
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.BuscarSeguindo(usuarioID)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao buscar seguidores: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, usuarios)
}

func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["id"], 10, 32)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o parâmetro: ", erro.Error())
		return
	}

	usuarioIDNoToken, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusUnauthorized, "Usuário não autorizado: ", erro.Error())
	}

	if usuarioID != usuarioIDNoToken {
		respostas.ResponderErro(w, http.StatusForbidden, "Operação não permitida: ", "Não é possível atualizar a senha de um usuário que não seja o seu")
		return
	}

	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao ler o body: ", erro.Error())
		return
	}

	var senha modelos.Senha

	if erro = json.Unmarshal(corpoRequest, &senha); erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o body: ", erro.Error())
		return
	}

	if senha.Nova == "" {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter o body: ", errors.New("É obrigatório informar a nova senha").Error())
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDeUsuarios(db)

	senhaSalvaNoBanco, erro := repositorio.BuscarSenha(usuarioID)

	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao buscar a senha no banco: ", erro.Error())
		return
	}

	if erro = seguranca.VerificarSenha(senhaSalvaNoBanco, senha.Atual); erro != nil {
		respostas.ResponderErro(w, http.StatusUnauthorized, "Senha atual não corresponde à salva no banco: ", erro.Error())
		return
	}

	senhaComHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao converter a nova senha: ", erro.Error())
		return
	}

	if erro = repositorio.AtualizarSenha(usuarioID, string(senhaComHash)); erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao atualizar a senha do usuário: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, "Usuário atualizado com sucesso!")
}
