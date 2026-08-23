package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"backend/src/autenticacao"
	"backend/src/banco"
	"backend/src/modelos"
	"backend/src/repositorios"
	"backend/src/respostas"

	"github.com/gorilla/mux"
)

func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao extrair ID do usuário: ", erro.Error())
		return
	}

	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao ler o body: ", erro.Error())
		return
	}

	var publicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequest, &publicacao); erro != nil {
		respostas.ResponderErro(w, http.StatusUnprocessableEntity, "Erro ao ler o json do body: ", erro.Error())
		return
	}

	if erro = publicacao.Preparar(); erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao processar o body: ", erro.Error())
		return
	}

	publicacao.AutorID = usuarioID

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	publicacao.ID, erro = repositorio.CriarPublicacao(publicacao)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao criar publicação: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, publicacao)
}

func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao extrair ID do usuário: ", erro.Error())
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao conectar ao banco: ", erro.Error())
		return
	}
	defer db.Close()

	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.BuscarPublicacoes(usuarioID)

	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao buscar publicacoes: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, publicacoes)
}

func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["id"], 10, 32)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacao, erro := repositorio.BuscarPublicacaoPorID(publicacaoID)

	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao consultar o banco: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, publicacao)
}

func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao extrair ID do usuário: ", erro.Error())
		return
	}

	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["id"], 10, 32)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	publicacao, erro := repositorio.BuscarPublicacaoPorID(publicacaoID)

	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao consultar o banco: ", erro.Error())
		return
	}

	if publicacao.AutorID != usuarioID {
		respostas.ResponderErro(w, http.StatusForbidden, "Erro ao atualizar a publicação: ", errors.New("Não é possível atualizar uma publicação que não seja sua").Error())
		return
	}

	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao ler o body: ", erro.Error())
		return
	}

	var novaPublicacao modelos.Publicacao
	if erro = json.Unmarshal(corpoRequest, &novaPublicacao); erro != nil {
		respostas.ResponderErro(w, http.StatusUnprocessableEntity, "Erro ao ler o json do body: ", erro.Error())
		return
	}

	if erro = novaPublicacao.Preparar(); erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao processar o body: ", erro.Error())
		return
	}

	novaPublicacao.ID = publicacao.ID

	if erro = repositorio.AtualizarPublicacao(novaPublicacao); erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao atualizar a publicação: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, "Publicação atualizada com sucesso.")
}

func ExcluirPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.ResponderErro(w, http.StatusBadRequest, "Erro ao extrair ID do usuário: ", erro.Error())
		return
	}

	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["id"], 10, 32)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	publicacao, erro := repositorio.BuscarPublicacaoPorID(publicacaoID)

	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao consultar o banco: ", erro.Error())
		return
	}

	if publicacao.AutorID != usuarioID {
		respostas.ResponderErro(w, http.StatusForbidden, "Erro ao excluir a publicação: ", errors.New("Não é possível excluir uma publicação que não seja sua").Error())
		return
	}

	if erro = repositorio.ExcluirPublicacao(publicacaoID); erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao excluir a publicação: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, "Sucesso ao excluir a publicação")
}

func BuscarPublicacoesUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(parametros["usuarioId"], 10, 32)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)
	publicacoes, erro := repositorio.BuscarPublicacoesUsuario(usuarioID)

	if erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao buscar publicacoes: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, publicacoes)
}

func CurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["id"], 10, 32)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	if erro = repositorio.CurtirPublicacao(publicacaoID); erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao curtir a publicação: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, "Sucesso ao curtir a publicação")
}

func DescurtirPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["id"], 10, 32)
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
	repositorio := repositorios.NovoRepositorioDePublicacoes(db)

	if erro = repositorio.DescurtirPublicacao(publicacaoID); erro != nil {
		respostas.ResponderErro(w, http.StatusInternalServerError, "Erro ao descurtir a publicação: ", erro.Error())
		return
	}

	respostas.ResponderJSON(w, http.StatusOK, "Sucesso ao descurtir a publicação")
}
