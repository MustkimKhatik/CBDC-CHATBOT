package utils

import (
    "os"
)

type Config struct {
    QdrantURL string
    OllamaURL string
    EmbedURL  string
}

func LoadConfig() Config {
    return Config{
        QdrantURL: getenvDefault("QDRANT_URL", "http://localhost:6333"),
        OllamaURL: getenvDefault("OLLAMA_URL", "http://localhost:11434"),
        EmbedURL:  getenvDefault("EMBED_URL", "http://localhost:8000/embed"),
    }
}

func getenvDefault(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}


