package utils

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type QdrantClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewQdrantClient(baseURL string) *QdrantClient {
    return &QdrantClient{baseURL: baseURL, httpClient: &http.Client{}}
}

type Point struct {
    ID      interface{}            `json:"id"`
    Vector  []float32              `json:"vector"`
    Payload map[string]interface{} `json:"payload"`
}

// UpsertPoints inserts or updates points into a collection.
func (q *QdrantClient) UpsertPoints(collection string, points []Point) error {
    body := map[string]interface{}{
        "points": points,
    }
    b, _ := json.Marshal(body)
    url := fmt.Sprintf("%s/collections/%s/points?wait=true", q.baseURL, collection)
    req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, err := q.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        msg, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("qdrant upsert status: %s - %s", resp.Status, string(msg))
    }
    return nil
}

type SearchRequest struct {
    Vector     []float32           `json:"vector"`
    Limit      int                 `json:"limit"`
    WithVector bool                `json:"with_vector"`
    WithPayload bool               `json:"with_payload"`
    ScoreThreshold *float32        `json:"score_threshold,omitempty"`
}

type SearchResult struct {
    Result []struct {
        ID      interface{}            `json:"id"`
        Score   float64                `json:"score"`
        Payload map[string]interface{} `json:"payload"`
    } `json:"result"`
}

func (q *QdrantClient) Search(collection string, reqBody SearchRequest) (SearchResult, error) {
    var result SearchResult
    b, _ := json.Marshal(reqBody)
    url := fmt.Sprintf("%s/collections/%s/points/search", q.baseURL, collection)
    req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, err := q.httpClient.Do(req)
    if err != nil {
        return result, err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return result, fmt.Errorf("qdrant search status: %s", resp.Status)
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return result, err
    }
    return result, nil
}

// CreateCollection ensures collection exists with correct vector size.
func (q *QdrantClient) CreateCollection(collection string, vectorSize int) error {
    body := map[string]interface{}{
        "vectors": map[string]interface{}{
            "size":     vectorSize,
            "distance": "Cosine",
        },
    }
    b, _ := json.Marshal(body)
    url := fmt.Sprintf("%s/collections/%s", q.baseURL, collection)
    req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(b))
    req.Header.Set("Content-Type", "application/json")
    resp, err := q.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode == http.StatusConflict { // 409 - already exists
        return nil
    }
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return fmt.Errorf("qdrant create collection status: %s", resp.Status)
    }
    return nil
}


