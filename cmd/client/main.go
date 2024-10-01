package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatalf("Erro: Timeout atingido na requisição. O servidor não respondeu dentro de 300ms.")
		} else {
			log.Fatalf("Erro ao fazer requisição: %v", err)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Erro: Servidor retornou status %d", resp.StatusCode)
		return
	}

	log.Println("Requisição feita com sucesso")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler corpo da resposta: %v", err)
	}

	writeFile(string(body))
}

func writeFile(cotacao string) {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("Erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("Dólar " + cotacao)
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}

	log.Println("Cotação salva com sucesso")
}
