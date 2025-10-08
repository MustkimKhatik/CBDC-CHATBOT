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

---

## New Machine Setup Guide

This guide shows two ways to run the stack: Docker Compose (recommended) or manual local run (Windows-friendly).

### Option A: Full stack with Docker Compose (recommended)
Prereqs: Docker Desktop (engine running)

1) Start all services
```
docker compose up --build -d
```

2) Pull llama3 model into Ollama container (once)
```
docker compose exec ollama ollama pull llama3
```

3) Open the apps
- Backend: `http://localhost:8080/api/health`
- Frontend: `http://localhost:5173`

4) Upload a document and ask a question
- Upload `.pdf/.txt/.md`
- Ask a question; answers are restricted to retrieved context

Notes
- PDF works out-of-the-box in Compose (backend image installs `poppler-utils`).
- Data is not persisted across container rebuilds by default except Qdrant volumes; add bind mounts if you need persistent uploads.

### Option B: Manual local run (Windows)
Prereqs: Go 1.22+, Node 20+, Python 3.11+

1) Start Qdrant
- Easiest via Docker (requires Docker Desktop):
```
docker run -d --name qdrant -p 6333:6333 -v qdrant_storage:/qdrant/storage qdrant/qdrant:latest
```
Or install Qdrant natively and run on port 6333.

2) Start Ollama (and pull llama3)
- Docker:
```
docker run -d --name ollama -p 11434:11434 -v ollama_models:/root/.ollama ollama/ollama:latest
docker exec -it ollama ollama pull llama3
```
Or install native Ollama, then run:
```
ollama serve
ollama pull llama3
```

3) Start the Embedding service
```
cd embed_service
python -m venv .venv
.\.venv\Scripts\activate
pip install -r requirements.txt
python app.py
```
It listens on `http://0.0.0.0:8000`. If you see IPv6 connection issues, the backend can be pointed at IPv4:
```
$env:EMBED_URL = 'http://127.0.0.1:8000/embed'
```

4) Start the Backend (Go)
```
cd backend
go mod tidy
go run .
```
Environment variables (override defaults as needed):
```
$env:QDRANT_URL = 'http://localhost:6333'
$env:OLLAMA_URL = 'http://localhost:11434'
$env:EMBED_URL  = 'http://127.0.0.1:8000/embed'
```

5) Start the Frontend (Vite)
```
cd frontend
npm install
npm run dev -- --host
```
Open `http://localhost:5173`.

### PDF Support (local Windows)
PDF extraction uses Poppler's `pdftotext` for robustness. You have two choices:
- Use Docker Compose (backend image already has `poppler-utils`) — no extra setup.
- Or install Poppler for Windows, then either add `pdftotext.exe` to PATH or set an explicit path for the backend process:
```
$env:PDFTOTEXT_PATH='C:\\Path\\to\\poppler\\bin\\pdftotext.exe'
go run .
```

If you prefer 100% pure-Go without Poppler, we can add a fallback extractor (lower fidelity). Open an issue.

### Health Checks
- Qdrant: `curl http://localhost:6333/readyz` → `all shards are ready`
- Ollama: `curl http://localhost:11434/api/tags` → lists models; ensure `llama3` appears
- Embed: `Invoke-WebRequest -Method POST -Uri http://127.0.0.1:8000/embed -Body (@{texts=@('hello')} | ConvertTo-Json) -ContentType 'application/json'`
- Backend: `curl http://localhost:8080/api/health`

### Common Errors & Fixes
- Embed 500 / numpy not available: ensure `numpy` is in `embed_service/requirements.txt` and rebuild the `embed` image.
- Backend cannot reach embed: start embedding service, then set `EMBED_URL` to IPv4 (`http://127.0.0.1:8000/embed`).
- Qdrant 400 vector size: delete the collection `cbdc_docs` and retry upload; backend recreates with 384 dims.
- Qdrant 409 create collection: safe to ignore; backend now treats it as non-error.
- Qdrant invalid point ID: fixed by using UUID v4 IDs.
- PDF `pdftotext` not found (Windows): install Poppler or set `PDFTOTEXT_PATH` env var.

### RAG Behavior Guarantees
- The backend declines to answer if no relevant context is retrieved.
- The prompt instructs the model to ONLY use provided context; if context lacks the answer, it must say it does not know based on documents.




