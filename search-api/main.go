package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"search-api/consumer"
	"search-api/solr"
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "El parámetro 'q' es requerido", http.StatusBadRequest)
		return
	}

	results, err := solr.SearchHotels(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al buscar hoteles: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	fmt.Println("Iniciando el consumidor de RabbitMQ...")
	go func() {
		err := consumer.Consume()
		if err != nil {
			log.Fatalf("Error al consumir los mensajes: %v", err)
		}
	}()

	// Nuevo endpoint para búsquedas
	http.HandleFunc("/search", searchHandler)

	fmt.Println("Servidor corriendo en :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
