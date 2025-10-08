package utils

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type EmbedClient struct {
    url        string
    httpClient *http.Client
}

func NewEmbedClient(url string) *EmbedClient {
    return &EmbedClient{url: url, httpClient: &http.Client{}}
}

type embedRequest struct {
    Texts []string `json:"texts"`
}

type embedResponse struct {
    Vectors [][]float32 `json:"vectors"`
}

func (e *EmbedClient) Embed(texts []string) ([][]float32, error) {
    payload := embedRequest{Texts: texts}
    b, _ := json.Marshal(payload)
    req, _ := http.NewRequest(http.MethodPost, e.url, bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, err := e.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("embed status: %s", resp.Status)
    }
    var out embedResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, err
    }
    return out.Vectors, nil
}


