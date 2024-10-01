package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/andrefelizardo/goposexpert_client-server/configs"
)

type Cotacao struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

var db *sql.DB

func main() {
	var err error
	db, err = configs.InitDB()
	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/cotacao", handlerCotacao)

	log.Println("Servidor iniciando na porta :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

func handlerCotacao(w http.ResponseWriter, r *http.Request) {
	cotacao, err := getCotacaoDolar()
	if err != nil {
		http.Error(w, "Erro ao obter cotação", http.StatusInternalServerError)
		return
	}

	err = armazenaCotacaoDolar(cotacao)
	if err != nil {
		http.Error(w, "Erro ao armazenar cotação", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cotacao))
}

func getCotacaoDolar() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Printf("Erro ao criar requisição: %v", err)
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Erro: Timeout atingido na requisição. O servidor não respondeu dentro de 200ms.")
		} else {
			log.Printf("Erro ao fazer requisição: %v", err)
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro: Servidor da API externa retornou status %d", resp.StatusCode)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler corpo da resposta: %v", err)
		return "", err
	}

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Printf("Erro ao fazer parse do JSON: %v", err)
		return "", err
	}

	log.Println("Cotação obtida com sucesso")
	return cotacao.Usdbrl.Bid, nil
}

func armazenaCotacaoDolar(cotacao string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	smt, err := db.PrepareContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)")
	if err != nil {
		log.Printf("Erro ao preparar statement: %v", err)
		return err
	}
	defer smt.Close()

	_, err = smt.ExecContext(ctx, cotacao)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Erro: Timeout atingido na execução do statement. O servidor não respondeu dentro de 10ms.")
		} else {
			log.Printf("Erro ao executar statement: %v", err)
		}
		return err
	}

	return nil
}
