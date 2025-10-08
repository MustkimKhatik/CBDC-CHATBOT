package utils

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type OllamaClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewOllamaClient(baseURL string) *OllamaClient {
    return &OllamaClient{baseURL: baseURL, httpClient: &http.Client{}}
}

type ollamaRequest struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
    Stream bool   `json:"stream"`
}

type ollamaResponse struct {
    Response string `json:"response"`
    Done     bool   `json:"done"`
}

func (o *OllamaClient) Generate(model, prompt string) (string, error) {
    payload := ollamaRequest{Model: model, Prompt: prompt, Stream: false}
    b, _ := json.Marshal(payload)
    req, _ := http.NewRequest(http.MethodPost, o.baseURL+"/api/generate", bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, err := o.httpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        // read body for better error reporting
        var buf bytes.Buffer
        _, _ = buf.ReadFrom(resp.Body)
        return "", fmt.Errorf("ollama status: %s - %s", resp.Status, buf.String())
    }
    var out ollamaResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return "", err
    }
    return out.Response, nil
}


