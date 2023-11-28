package treblle_fiber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	timeoutDuration = 2 * time.Second
)

func sendToTreblle(treblleInfo MetaData) {
	bytesRepresentation, err := json.Marshal(treblleInfo)
	if err != nil {
		fmt.Printf("failed to marshall treblle information: %+v\n", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, Config.ServerURL, bytes.NewBuffer(bytesRepresentation))
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
			bytes := []byte{}
			if _, err = response.Body.Read(bytes); err != nil {
				fmt.Printf("failed to post: %+v\n", err)
			}
			fmt.Printf("failed to post: %+v\n", string(bytes))
		}
	}
}
