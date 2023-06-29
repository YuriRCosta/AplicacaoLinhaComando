package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const ERRO_AO_CONVERTER_O_PARAMETRO_PARA_INTEIRO = "Erro ao converter o parâmetro para inteiro!"
const ERRO_AO_CRIAR_STATEMENT = "Erro ao criar o statement!"
const ERRO_AO_CONECTAR_COM_O_BANCO_DE_DADOS = "Erro ao conectar com o banco de dados!"

type usuario struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

// CriarUsuario insere um usuário no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição!"))
		return
	}

	var usuario usuario

	if err = json.Unmarshal(corpoRequisicao, &usuario); err != nil {
		w.Write([]byte("Erro ao converter o usuário para struct!"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte(ERRO_AO_CONECTAR_COM_O_BANCO_DE_DADOS))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("insert into usuarios (nome, email) values (?, ?)")
	if err != nil {
		w.Write([]byte(ERRO_AO_CRIAR_STATEMENT))
		return
	}
	defer statement.Close()

	insercao, err := statement.Exec(usuario.Nome, usuario.Email)
	if err != nil {
		w.Write([]byte("Erro ao executar o statement!"))
		return
	}

	idInserido, err := insercao.LastInsertId()
	if err != nil {
		w.Write([]byte("Erro ao obter o ID inserido!"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Usuário inserido com sucesso! ID: %d", idInserido)))

	fmt.Println(usuario)
}

// BuscarUsuarios busca todos os usuários salvos no banco de dados
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte(ERRO_AO_CONECTAR_COM_O_BANCO_DE_DADOS))
		return
	}
	defer db.Close()

	linhas, err := db.Query("select * from usuarios")
	if err != nil {
		w.Write([]byte("Erro ao buscar os usuários!"))
		return
	}
	defer linhas.Close()

	var usuarios []usuario

	for linhas.Next() {
		var usuario usuario

		if err := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("Erro ao escanear o usuário!"))
			return
		}

		usuarios = append(usuarios, usuario)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		w.Write([]byte("Erro ao converter os usuários para JSON!"))
		return
	}
}

// BuscarUsuario busca um usuário salvo no banco de dados
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte(ERRO_AO_CONVERTER_O_PARAMETRO_PARA_INTEIRO))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte(ERRO_AO_CONECTAR_COM_O_BANCO_DE_DADOS))
		return
	}
	defer db.Close()

	linha, err := db.Query("select * from usuarios where id = ?", ID)
	if err != nil {
		w.Write([]byte("Erro ao buscar o usuário!"))
		return
	}
	defer linha.Close()

	var usuario usuario
	if linha.Next() {
		if err := linha.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("Erro ao escanear o usuário!"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usuario); err != nil {
		w.Write([]byte("Erro ao converter o usuário para JSON!"))
		return
	}
}

// AtualizarUsuario altera os dados de um usuário no banco de dados
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte(ERRO_AO_CONVERTER_O_PARAMETRO_PARA_INTEIRO))
		return
	}

	corpoRequisicao, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição!"))
		return
	}

	var usuario usuario

	if err = json.Unmarshal(corpoRequisicao, &usuario); err != nil {
		w.Write([]byte("Erro ao converter o usuário para struct!"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte(ERRO_AO_CONECTAR_COM_O_BANCO_DE_DADOS))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update usuarios set nome = ?, email = ? where id = ?")
	if err != nil {
		w.Write([]byte(ERRO_AO_CRIAR_STATEMENT))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(usuario.Nome, usuario.Email, ID); err != nil {
		w.Write([]byte("Erro ao atualizar o usuário!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeletarUsuario remove um usuário do banco de dados
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		w.Write([]byte(ERRO_AO_CONVERTER_O_PARAMETRO_PARA_INTEIRO))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte(ERRO_AO_CONECTAR_COM_O_BANCO_DE_DADOS))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("delete from usuarios where id = ?")
	if err != nil {
		w.Write([]byte(ERRO_AO_CRIAR_STATEMENT))
		return
	}
	defer statement.Close()

	if _, err := statement.Exec(ID); err != nil {
		w.Write([]byte("Erro ao deletar o usuário!"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
