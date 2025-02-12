package solr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Estructura de ejemplo para los documentos de hoteles
type Hotel struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Country string `json:"country"`
}

// Enviar el hotel a Solr
func SendToSolr(message []byte) error {
	var hotel Hotel
	err := json.Unmarshal(message, &hotel)
	if err != nil {
		return fmt.Errorf("error unmarshaling message: %v", err)
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"add": map[string]interface{}{
			"doc": hotel,
		},
	})
	if err != nil {
		return fmt.Errorf("error marshaling data to Solr: %v", err)
	}

	solrURL := "http://localhost:8983/solr/hotel_core/update?commit=true"
	req, err := http.NewRequest("POST", solrURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Solr: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("solr returned an error: %v", resp.Status)
	}

	log.Println("Document successfully indexed in Solr.")
	return nil
}

// SearchHotels busca hoteles en Solr según una consulta de búsqueda
func SearchHotels(queryTerm string) (map[string]interface{}, error) {
	// Construcción de la consulta para búsqueda global
	query := fmt.Sprintf("%s*", queryTerm)
	solrURL := fmt.Sprintf("http://localhost:8983/solr/hotel_core/select?q=%s&wt=json", query)

	resp, err := http.Get(solrURL)
	if err != nil {
		return nil, fmt.Errorf("error al hacer la solicitud a Solr: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error al decodificar la respuesta de Solr: %v", err)
	}

	return result, nil
}
