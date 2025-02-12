package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

// URL de tu Solr
const solrURL = "http://localhost:8983/solr/hotel_core/update?commit=true"

// Consume escucha los mensajes de la cola "hotel_created"
func Consume() error {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Error al conectar a RabbitMQ: %v", err)
		return err
	}
	defer conn.Close()

	// Crear un canal de comunicación
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error al crear un canal: %v", err)
		return err
	}
	defer ch.Close()

	// Declarar la cola
	queue, err := ch.QueueDeclare(
		"hotel_created", // Nombre de la cola
		true,            // durable
		false,           // auto-deleted
		false,           // exclusive
		false,           // no-wait
		nil,             // argumentos
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %v", err)
		return err
	}

	// Consumir los mensajes de la cola
	msgs, err := ch.Consume(
		queue.Name, // nombre de la cola
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // argumentos
	)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
		return err
	}

	// Procesar los mensajes
	go func() {
		for msg := range msgs {
			fmt.Printf("Mensaje recibido: %s\n", msg.Body)

			// Convertir el mensaje a un objeto hotel
			var hotel map[string]interface{}
			err := json.Unmarshal(msg.Body, &hotel)
			if err != nil {
				log.Printf("Error al deserializar el mensaje: %v", err)
				continue
			}

			// Indexar el hotel en Solr
			err = indexHotelInSolr(hotel)
			if err != nil {
				log.Printf("Error al indexar el hotel en Solr: %v", err)
			}
		}
	}()

	// Esperar indefinidamente para seguir recibiendo mensajes
	select {}
}

// indexHotelInSolr indexa un hotel en Solr
func indexHotelInSolr(hotel map[string]interface{}) error {
	// Crear un documento Solr con el hotel
	log.Println("indexando hotel en solr", hotel)

	solrDoc := map[string]interface{}{
		"add": map[string]interface{}{
			"doc": hotel,
		},
	}

	// Convertir el documento Solr a JSON
	hotelJSON, err := json.Marshal(solrDoc)
	if err != nil {
		log.Println("Error al convertir hotel a JSON:", err)
		return err
	}

	// Preparar la solicitud HTTP
	req, err := http.NewRequest("POST", solrURL, bytes.NewBuffer(hotelJSON))
	if err != nil {
		log.Println("Error al crear solicitud HTTP:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud a Solr
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error al enviar solicitud a Solr:", err)
		return err
	}
	defer resp.Body.Close()

	// Verificar la respuesta
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error al indexar en Solr. Código de respuesta: %d\n", resp.StatusCode)
		return fmt.Errorf("error al indexar en Solr. Código de respuesta: %d", resp.StatusCode)
	}

	log.Println("Hotel indexado en Solr con éxito")
	return nil
}
