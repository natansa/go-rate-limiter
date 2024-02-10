package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	var wg sync.WaitGroup
	requestsPerSecond := 1 // quantidade de requisições por segundo
	totalRequests := 20    // quantidade total de requisições
	apiKey := "API_KEY_A"  // API_KEY_A | API_KEY_B | API_KEY_C

	// WaitGroup para esperar até que todas as goroutines sejam concluídas
	wg.Add(totalRequests)

	// função anônima para representar cada requisição
	makeRequest := func() {
		defer wg.Done()

		// nova solicitação HTTP
		req, err := http.NewRequest("GET", "http://localhost:8080", nil)
		if err != nil {
			fmt.Printf("Erro ao criar a requisição: %v\n", err)
			return
		}

		// adiciona a API_KEY ao cabeçalho da requisição
		req.Header.Set("API_KEY", apiKey)

		// envia a requsiçao
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Erro na requisição: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			// le o body do response
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// retorna a mensagem do response
			message := string(body)
			fmt.Printf("Result: %s\n", message)
		}
	}

	// inicia as goroutines para simular as requisições
	for i := 0; i < totalRequests; i++ {
		go makeRequest()
		time.Sleep(time.Second / time.Duration(requestsPerSecond))
	}

	// aguarda até que todas as goroutines sejam concluídas
	wg.Wait()
}
