package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PerencanaanClientImpl struct {
	BaseClient
}

func NewPerencanaanClient(httpClient *http.Client) *PerencanaanClientImpl {
	return &PerencanaanClientImpl{
		BaseClient: newBaseClient(
			"https://api-perencanaan-dev-mahulu.zeabur.app",
			"",
			httpClient,
		),
	}
}

func (c *PerencanaanClientImpl) GetPotensiPerangkatDaerah(ctx context.Context, ids []int) ([]PpdItem, error) {
	// url check program unggulan
	url := fmt.Sprintf("%s/ppd/find-by-ids", c.host)
	// body kode program unggulans
	payload := FindByIdsRequest{
		Id: ids,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("gagal encode body: %w", err)
	}

	// request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	sessionID := getSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Tidak ada Session Id ditemukan, mungkin akan 401")
	}

	// send request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request ke id gagal: %w", err)
	}
	defer res.Body.Close()

	// response status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Id: %v tidak ditemukan. status: %d", ids, res.StatusCode)
	}

	var result PotensiPerangkatDaerahResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response: %w", err)
	}
	return result.Data, nil

}
func (c *PerencanaanClientImpl) GetIsuKlhs(ctx context.Context, ids []int) ([]IsuItem, error) {
	// url check program unggulan
	url := fmt.Sprintf("%s/isu-klhs/find-by-ids", c.host)
	// body kode program unggulans
	payload := FindByIdsRequest{
		Id: ids,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("gagal encode body: %w", err)
	}

	// request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	sessionID := getSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Tidak ada Session Id ditemukan, mungkin akan 401")
	}

	// send request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request ke id gagal: %w", err)
	}
	defer res.Body.Close()

	// response status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Id: %v tidak ditemukan. status: %d", ids, res.StatusCode)
	}

	var result IsuResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response: %w", err)
	}
	return result.Data, nil

}
func (c *PerencanaanClientImpl) GetIsuGlobal(ctx context.Context, ids []int) ([]IsuItem, error) {
	// url check program unggulan
	url := fmt.Sprintf("%s/isu-global/find-by-ids", c.host)
	// body kode program unggulans
	payload := FindByIdsRequest{
		Id: ids,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("gagal encode body: %w", err)
	}

	// request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	sessionID := getSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Tidak ada Session Id ditemukan, mungkin akan 401")
	}

	// send request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request ke id gagal: %w", err)
	}
	defer res.Body.Close()

	// response status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Id: %v tidak ditemukan. status: %d", ids, res.StatusCode)
	}

	var result IsuResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response: %w", err)
	}
	return result.Data, nil

}
func (c *PerencanaanClientImpl) GetIsuNasional(ctx context.Context, ids []int) ([]IsuItem, error) {
	// url check program unggulan
	url := fmt.Sprintf("%s/isu-nasional/find-by-ids", c.host)
	// body kode program unggulans
	payload := FindByIdsRequest{
		Id: ids,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("gagal encode body: %w", err)
	}

	// request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	sessionID := getSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Tidak ada Session Id ditemukan, mungkin akan 401")
	}

	// send request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request ke id gagal: %w", err)
	}
	defer res.Body.Close()

	// response status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Id: %v tidak ditemukan. status: %d", ids, res.StatusCode)
	}

	var result IsuResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response: %w", err)
	}
	return result.Data, nil

}
func (c *PerencanaanClientImpl) GetIsuRegional(ctx context.Context, ids []int) ([]IsuItem, error) {
	// url check program unggulan
	url := fmt.Sprintf("%s/isu-regional/find-by-ids", c.host)
	// body kode program unggulans
	payload := FindByIdsRequest{
		Id: ids,
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("gagal encode body: %w", err)
	}

	// request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("Gagal membuat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	sessionID := getSessionID(ctx)
	if sessionID != "" {
		req.Header.Set("X-Session-Id", sessionID)
	} else {
		log.Printf("Tidak ada Session Id ditemukan, mungkin akan 401")
	}

	// send request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request ke id gagal: %w", err)
	}
	defer res.Body.Close()

	// response status
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Id: %v tidak ditemukan. status: %d", ids, res.StatusCode)
	}

	var result IsuResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gagal decode response: %w", err)
	}
	return result.Data, nil

}