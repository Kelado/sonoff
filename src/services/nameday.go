package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type NamedayService struct {
}

func NewNamedayService() *NamedayService {
	return &NamedayService{}
}

func (n *NamedayService) GetCelebratingNamesForToday() []string {
	return getCelebratingNamesForToday()
}

func getCelebratingNamesForToday() []string {
	url := "https://nameday.abalin.net/api/V2/today?timezone=Europe/Athens"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var response struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return nil
	}

	names, ok := response.Data["gr"]
	if !ok {
		fmt.Println("No names found for 'gr'")
		return nil
	}

	return splitNames(names)
}

func splitNames(names string) []string {
	return strings.Split(names, ", ")
}
