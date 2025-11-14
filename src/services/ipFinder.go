package services

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
)

const ApiURL = "https://api.ipify.org/"

type IpService struct {
}

func NewIpService() *IpService {
	return &IpService{}
}

func (s *IpService) GetPublicIp() (net.IP, error) {
	req, err := http.NewRequest(http.MethodGet, ApiURL, nil)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if err == nil {
			return nil, fmt.Errorf(
				"api responded with status code %d, message '%s'",
				res.StatusCode,
				string(body),
			)
		} else {
			return nil, fmt.Errorf("api responded with status code %d", res.StatusCode)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	bodyString := string(body)
	ip := net.ParseIP(bodyString)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse api response '%s' as IP", bodyString)
	}

	return ip, nil
}
