package utils

import (
    "strings"
)

// SplitTextIntoChunks splits text into roughly sized chunks by sentences and paragraphs.
// targetTokens is approximate, using word count as a proxy.
func SplitTextIntoChunks(text string, targetTokens int) []string {
    if targetTokens <= 0 {
        targetTokens = 1000
    }

    // Normalize newlines
    text = strings.ReplaceAll(text, "\r\n", "\n")
    text = strings.TrimSpace(text)
    if text == "" {
        return nil
    }

    // Simple heuristic: split by double newline paragraphs, then sentences.
    paragraphs := strings.Split(text, "\n\n")
    var chunks []string
    var current []string
    currentWords := 0

    flush := func() {
        if len(current) == 0 {
            return
        }
        chunks = append(chunks, strings.TrimSpace(strings.Join(current, " ")))
        current = current[:0]
        currentWords = 0
    }

    for _, p := range paragraphs {
        p = strings.TrimSpace(p)
        if p == "" {
            continue
        }
        sentences := splitIntoSentences(p)
        for _, s := range sentences {
            w := wordCount(s)
            if currentWords+w > targetTokens && currentWords > 0 {
                flush()
            }
            current = append(current, s)
            currentWords += w
        }
    }
    flush()
    return chunks
}

func wordCount(s string) int {
    s = strings.TrimSpace(s)
    if s == "" {
        return 0
    }
    return len(strings.Fields(s))
}

// splitIntoSentences is a naive sentence splitter suitable for tech text.
func splitIntoSentences(p string) []string {
    // naive split on period/question/exclamation followed by space/newline
    // keep delimiters
    var out []string
    start := 0
    for i := 0; i < len(p); i++ {
        c := p[i]
        if c == '.' || c == '!' || c == '?' {
            // include delimiter
            j := i + 1
            // also handle multiple punctuation like ") ."
            for j < len(p) && (p[j] == ' ' || p[j] == '\n') {
                j++
                break
            }
            out = append(out, strings.TrimSpace(p[start:j]))
            start = j
        }
    }
    if start < len(p) {
        out = append(out, strings.TrimSpace(p[start:]))
    }
    if len(out) == 0 {
        return []string{strings.TrimSpace(p)}
    }
    return out
}


