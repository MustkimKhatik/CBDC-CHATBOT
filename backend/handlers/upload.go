package handlers

import (
    "fmt"
    "io"
    "net/http"
    "path/filepath"
    "strings"

    "github.com/gin-gonic/gin"

    "cbdc-backend/utils"
    "github.com/google/uuid"
)

type uploadResponse struct {
    FileName   string `json:"file_name"`
    NumChunks  int    `json:"num_chunks"`
    Collection string `json:"collection"`
}

// Upload handles document upload, extraction, chunking, embedding, and Qdrant upsert.
func Upload(c *gin.Context) {
    cfg := utils.LoadConfig()
    qdrant := utils.NewQdrantClient(cfg.QdrantURL)
    embed := utils.NewEmbedClient(cfg.EmbedURL)

    // Ensure collection exists (vector size 384 for all-MiniLM-L6-v2)
    if err := qdrant.CreateCollection("cbdc_docs", 384); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("create collection: %v", err)})
        return
    }

    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
        return
    }
    defer file.Close()

    ext := strings.ToLower(filepath.Ext(header.Filename))
    var text string
    switch ext {
    case ".pdf":
        text, err = utils.ExtractPDFText(file)
    case ".md", ".txt":
        b, rerr := io.ReadAll(file)
        if rerr != nil {
            err = rerr
            break
        }
        text = string(b)
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type"})
        return
    }
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("extract text: %v", err)})
        return
    }

    chunks := utils.SplitTextIntoChunks(text, 1000)
    if len(chunks) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no text found in document"})
        return
    }

    // Embed chunks
    vectors, err := embed.Embed(chunks)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("embed: %v", err)})
        return
    }
    if len(vectors) != len(chunks) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "embedding size mismatch"})
        return
    }

    // Upsert into Qdrant
    points := make([]utils.Point, 0, len(chunks))
    for i := range chunks {
        points = append(points, utils.Point{
            ID:     uuid.NewString(),
            Vector: vectors[i],
            Payload: map[string]interface{}{
                "text":     chunks[i],
                "file_name": header.Filename,
                "index":    i,
            },
        })
    }

    if err := qdrant.UpsertPoints("cbdc_docs", points); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("qdrant upsert: %v", err)})
        return
    }

    c.JSON(http.StatusOK, uploadResponse{FileName: header.Filename, NumChunks: len(chunks), Collection: "cbdc_docs"})
}


