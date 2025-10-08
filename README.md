## CBDC Assistant (Offline RAG)

Offline AI assistant for NPCI CBDC docs using Go (Gin), Qdrant, Ollama, and a local embedding microservice.

### Stack
- Backend: Go + Gin
- Vector DB: Qdrant
- LLM: Ollama (`llama3`)
- Embeddings: `all-MiniLM-L6-v2` (Flask microservice)
- Frontend: React + Vite
- Orchestration: Docker Compose

### Local run (dev)
Prereqs: Go 1.22+, Node 20+, Python 3.11+, Docker

Backend:
```
cd backend
go mod tidy
go run .
```

Frontend:
```
cd frontend
npm install
npm run dev
```

Embedding service (optional local run):
```
cd embed_service
pip install -r requirements.txt
python app.py
```

Health check:
```
GET http://localhost:8080/api/health
```

### Docker Compose
```
docker compose up --build
```
Services:
- Backend: http://localhost:8080
- Frontend: http://localhost:5173
- Qdrant: http://localhost:6333
- Ollama: http://localhost:11434
- Embed: http://localhost:8000

First run model pull (if needed):
```
docker exec -it $(docker ps -qf name=ollama) ollama pull llama3
```

### Env
Set environment variables (or use defaults):
- `QDRANT_URL` (default `http://localhost:6333`)
- `OLLAMA_URL` (default `http://localhost:11434`)
- `EMBED_URL` (default `http://localhost:8000/embed`)

### API
- `POST /api/upload` (form-data: `file` .pdf/.txt/.md)
- `POST /api/query`  JSON `{ "query": "..." }`
- `GET /api/health`



