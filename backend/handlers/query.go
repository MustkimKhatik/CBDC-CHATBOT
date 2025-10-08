package handlers

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"

    "cbdc-backend/utils"
)

type queryRequest struct {
    Query string `json:"query" binding:"required"`
}

type queryResponse struct {
    Answer   string   `json:"answer"`
    Contexts []string `json:"contexts"`
}

// Query performs RAG over Qdrant and generates an answer with Ollama.
func Query(c *gin.Context) {
    var req queryRequest
    if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Query) == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
        return
    }

    cfg := utils.LoadConfig()
    embed := utils.NewEmbedClient(cfg.EmbedURL)
    qdrant := utils.NewQdrantClient(cfg.QdrantURL)
    ollama := utils.NewOllamaClient(cfg.OllamaURL)

    // Embed query
    vectors, err := embed.Embed([]string{req.Query})
    if err != nil || len(vectors) == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("embed query: %v", err)})
        return
    }

    // Retrieve similar chunks
    searchRes, err := qdrant.Search("cbdc_docs", utils.SearchRequest{
        Vector:      vectors[0],
        Limit:       5,
        WithVector:  false,
        WithPayload: true,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("search: %v", err)})
        return
    }

    var contexts []string
    for _, r := range searchRes.Result {
        if t, ok := r.Payload["text"].(string); ok {
            contexts = append(contexts, t)
        }
    }
    if len(contexts) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no relevant context found; upload documents first or refine your query"})
        return
    }
    contextText := strings.Join(contexts, "\n---\n")

    // Build prompt
    prompt := fmt.Sprintf("SYSTEM: You are CBDC Assistant for NPCI. You must ONLY use the provided CONTEXT. If the answer is not fully contained in the CONTEXT, say 'I don't know based on the provided documents.' Never use prior knowledge.\nCONTEXT:\n%s\n---\nQUESTION: %s\nANSWER:", contextText, req.Query)

    // Generate answer
    answer, err := ollama.Generate("llama3", prompt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("ollama: %v", err)})
        return
    }

    c.JSON(http.StatusOK, queryResponse{Answer: strings.TrimSpace(answer), Contexts: contexts})
}


