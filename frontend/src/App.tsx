import React, { useState } from 'react'
import axios from 'axios'

export default function App() {
  const [file, setFile] = useState<File | null>(null)
  const [query, setQuery] = useState('')
  const [answer, setAnswer] = useState('')
  const [contexts, setContexts] = useState<string[]>([])
  const [busy, setBusy] = useState(false)
  const [drag, setDrag] = useState(false)
  const backendBase = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'

  const onUpload = async () => {
    if (!file) return
    setBusy(true)
    try {
      const form = new FormData()
      form.append('file', file)
      const res = await axios.post(`${backendBase}/api/upload`, form, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      alert(`Uploaded and indexed ${res.data.file_name} with ${res.data.num_chunks} chunks`)
    } catch (e: any) {
      const msg = e?.response?.data?.error || e?.message || 'Upload failed'
      alert(`Upload error: ${msg}`)
    } finally {
      setBusy(false)
    }
  }

  const onQuery = async () => {
    setAnswer('')
    setContexts([])
    try {
      const res = await axios.post(`${backendBase}/api/query`, { query })
      setAnswer(res.data.answer)
      setContexts(res.data.contexts || [])
    } catch (e: any) {
      const msg = e?.response?.data?.error || e?.message || 'Request failed'
      setAnswer(`Error: ${msg}`)
    }
  }

  const onDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setDrag(false)
    const f = e.dataTransfer.files?.[0]
    if (f) setFile(f)
  }

  return (
    <div className="container">
      <div className="header">
        <div className="title">CBDC Assistant</div>
        <div className="badge">Offline RAG · NPCI</div>
      </div>

      <div className="panel">
        <div className="section-title">Upload document</div>
        <div
          className={"dropzone" + (drag ? " drag" : "")}
          onDragOver={(e) => { e.preventDefault(); setDrag(true) }}
          onDragLeave={() => setDrag(false)}
          onDrop={onDrop}
        >
          {file ? (
            <div className="row">
              <span className="pill">{file.name}</span>
              <button className="btn" onClick={() => setFile(null)}>Clear</button>
            </div>
          ) : (
            <>
              Drag & drop .pdf/.txt/.md here or
              <label className="btn ghost" style={{ marginLeft: 8 }}>
                Choose file
                <input type="file" accept=".pdf,.txt,.md" onChange={(e) => setFile(e.target.files?.[0] || null)} style={{ display: 'none' }} />
              </label>
            </>
          )}
        </div>
        <div className="toolbar" style={{ marginTop: 10 }}>
          <button className="btn" onClick={() => setFile(null)} disabled={!file || busy}>Reset</button>
          <button className="btn primary" onClick={onUpload} disabled={!file || busy}>{busy ? 'Uploading…' : 'Upload & Index'}</button>
        </div>
      </div>

      <div className="panel">
        <div className="section-title">Ask</div>
        <div className="row" style={{ alignItems: 'stretch' }}>
          <input className="input" value={query} onChange={(e) => setQuery(e.target.value)} placeholder="How does inter-bank settlement work?" />
          <button className="btn secondary" onClick={onQuery}>Ask</button>
        </div>
        {answer && (
          <div style={{ marginTop: 14 }}>
            <div className="section-title">Answer</div>
            <div className="answer">{answer}</div>
            {contexts.length > 0 && (
              <div className="contexts">
                <div className="section-title">Context</div>
                <div className="accordion">
                  {contexts.map((c, i) => (
                    <details key={i} className="accordion-item">
                      <summary className="accordion-header">Chunk #{i + 1}</summary>
                      <pre className="accordion-content">{c}</pre>
                    </details>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}


