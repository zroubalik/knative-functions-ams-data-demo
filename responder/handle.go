package function

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudevents/sdk-go/v2/event"
)

func Handle(ctx context.Context, e event.Event) error {

	fmt.Println("Received event")
	fmt.Println(e)

	backendApiUri := os.Getenv("REACT_APP_BACKEND_URI")
	if backendApiUri == "" {
		backendApiUri = "http://ams-data-app-service:3333"
	}

	switch e.Type() {
	case "weather":
		err := sendDataToAPI(backendApiUri+"/weather", e.Data())
		if err != nil {
			return err
		}
	case "scooters":
		err := sendDataToAPI(backendApiUri+"/scooters", e.Data())
		if err != nil {
			return err
		}
	default:
		fmt.Println("Unknown event type:", e.Type())
	}

	return nil
}

func sendDataToAPI(url string, data []byte) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send data to API, status code: %d", resp.StatusCode)
	}
	fmt.Printf("Data were sent to %s, data: %s\n", url, string(data))
	return nil
}
