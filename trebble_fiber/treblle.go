package treblle_fiber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	timeoutDuration = 5 * time.Second
)

type BaseUrlOptions struct {
	Debug bool
}

func getTreblleBaseUrl() string {
	if Config.Debug {
		return "https://debug.treblle.com/"
	}

	treblleBaseUrls := []string{
		"https://rocknrolla.treblle.com",
		"https://punisher.treblle.com",
		"https://sicario.treblle.com",
	}

	rand.Seed(time.Now().Unix())
	randomUrlIndex := rand.Intn(len(treblleBaseUrls))

	return treblleBaseUrls[randomUrlIndex]
}

func sendToTreblle(treblleInfo MetaData) {
	baseUrl := getTreblleBaseUrl()

	bytesRepresentation, err := json.Marshal(treblleInfo)
	if err != nil {
		fmt.Printf("failed to marshall treblle information: %+v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, baseUrl, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Printf("failed to create HTTP Post request: %+v\n", err)
		return
	}
	// Set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", Config.APIKey)

	// Do the request
	client := &http.Client{Timeout: timeoutDuration}
	if response, err := client.Do(req); err != nil {
		fmt.Printf("failed to post: %+v\n", err)
	} else {
		if response.StatusCode >= http.StatusBadRequest {
			fmt.Printf("failed to post with status code >= 400: %+v\n", response)
		}
	}
}
