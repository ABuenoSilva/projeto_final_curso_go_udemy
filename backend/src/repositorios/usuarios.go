package repositorios

import (
	"database/sql"
	"errors"
	"fmt"

	"backend/src/modelos"
)

type usuarios struct {
	db *sql.DB
}

var ErrNenhumRegistroEncontrado = errors.New("Nenhum registro foi afetado!")

func NovoRepositorioDeUsuarios(db *sql.DB) *usuarios {
	return &usuarios{db}
}

func (repositorio usuarios) Criar(usuario modelos.Usuario) (uint64, error) {
	statement, erro := repositorio.db.Prepare("insert into usuarios (nome, nick, email, senha) values(?, ?, ?, ?)")
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, usuario.Senha)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

func (repositorio usuarios) Buscar(nomeOuNick string) ([]modelos.Usuario, error) {
	nomeOuNick = fmt.Sprintf("%%%s%%", nomeOuNick) // %<nomeOuNick>%
	linhas, erro := repositorio.db.Query("select id, nome, nick, email, criadoEm from usuarios where nome like ? or nick like ?", nomeOuNick, nomeOuNick)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	if erro = linhas.Err(); erro != nil {
		return nil, erro
	}

	var usuarios []modelos.Usuario

	for linhas.Next() {
		var usuario modelos.Usuario
		if erro = linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Nick, &usuario.Email, &usuario.CriadoEm); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
	}

	return usuarios, nil
}

func (repositorio usuarios) BuscarPorID(ID uint64) (modelos.Usuario, error) {
	linha, erro := repositorio.db.Query("select id, nome, nick, email, criadoEm from usuarios where id = ?", ID)
	if erro != nil {
		return modelos.Usuario{}, erro
	}
	defer linha.Close()

	if erro = linha.Err(); erro != nil {
		return modelos.Usuario{}, erro
	}

	var usuario modelos.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Nick, &usuario.Email, &usuario.CriadoEm); erro != nil {
			return modelos.Usuario{}, erro
		}
	}

	return usuario, nil
}

func (repositorio usuarios) Atualizar(ID uint64, usuario modelos.Usuario) error {
	statement, erro := repositorio.db.Prepare("update usuarios set nome = ?, nick = ?, email = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(usuario.Nome, usuario.Nick, usuario.Email, ID)
	if erro != nil {
		return erro
	}

	linhasAfetadas, erro := resultado.RowsAffected()
	if erro != nil {
		return erro
	}

	if linhasAfetadas == 0 {
		return ErrNenhumRegistroEncontrado
	}

	return nil
}

func (repositorio usuarios) Excluir(ID uint64) error {
	statement, erro := repositorio.db.Prepare("delete from usuarios where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(ID)

	if erro != nil {
		return erro
	}

	linhasAfetadas, erro := resultado.RowsAffected()
	if erro != nil {
		return erro
	}

	if linhasAfetadas == 0 {
		return ErrNenhumRegistroEncontrado
	}

	return nil
}
