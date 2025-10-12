package helper

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func FormatFile(fileBytes []byte, fileHeader *multipart.FileHeader) string {
	mimeType := http.DetectContentType(fileBytes)
	var formatExtension string
	switch mimeType {
	case "image/jpeg":
		formatExtension = ".jpg"
	case "image/png":
		formatExtension = ".png"
	case "application/pdf":
		formatExtension = ".pdf"
	case "video/mp4":
		formatExtension = ".mp4"
	// Tambahkan case untuk tipe file lain yang Anda dukung
	default:
		// Jika deteksi MIME tidak spesifik, gunakan ekstensi dari nama file
		formatExtension = strings.ToLower(filepath.Ext(fileHeader.Filename))
		// Jika ekstensi dari nama file juga tidak ada, gunakan hasil dari MIME
		if formatExtension == "" {
			formatExtension = "." + strings.Split(mimeType, "/")[1]
		}
	}

	return formatExtension

}

func FormatDateTimeToString(t time.Time) string {
	const datetimeFormat = "2006-01-02 15:04:05"
	return t.Format(datetimeFormat)
}
