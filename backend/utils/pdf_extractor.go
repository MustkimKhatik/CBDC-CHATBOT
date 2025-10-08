package utils

import (
    "bytes"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "time"
)

// ExtractPDFText uses Poppler's pdftotext to extract text content.
func ExtractPDFText(r io.Reader) (string, error) {
    // Write input to temp PDF
    tmpDir := os.TempDir()
    base := "cbdc_pdf_" + time.Now().Format("20060102150405")
    pdfPath := filepath.Join(tmpDir, base+".pdf")
    txtPath := filepath.Join(tmpDir, base+".txt")

    f, err := os.Create(pdfPath)
    if err != nil {
        return "", err
    }
    if _, err := io.Copy(f, r); err != nil {
        f.Close()
        os.Remove(pdfPath)
        return "", err
    }
    f.Close()
    defer os.Remove(pdfPath)
    defer os.Remove(txtPath)

    // Resolve pdftotext binary (supports Windows .exe and optional PDFTOTEXT_PATH)
    bin := os.Getenv("PDFTOTEXT_PATH")
    if bin == "" {
        if runtime.GOOS == "windows" {
            bin = "pdftotext.exe"
        } else {
            bin = "pdftotext"
        }
    }
    if _, err := exec.LookPath(bin); err != nil {
        return "", err
    }

    // Run: pdftotext -layout pdfPath txtPath
    cmd := exec.Command(bin, "-layout", pdfPath, txtPath)
    if err := cmd.Run(); err != nil {
        return "", err
    }

    b, err := os.ReadFile(txtPath)
    if err != nil {
        return "", err
    }
    return string(bytes.TrimSpace(b)), nil
}


