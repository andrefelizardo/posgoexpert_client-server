package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/cotacao", handlerCotacao)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
	log.Println("Servidor iniciado na porta :8080")

	
}

func handlerCotacao(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Cotação do dólar"}`))
}