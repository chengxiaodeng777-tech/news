package system

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestArticleImageUploadAndRead(t *testing.T) {
	t.Setenv("ARTICLE_UPLOAD_DIR", t.TempDir())
	router := gin.New()
	router.POST("/api/upload", UploadArticleImage)
	router.GET("/api/uploads/articles/:name", ServeArticleImage)

	png := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 0}
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	file, err := writer.CreateFormFile("file", "image.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write(png); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	result := httptest.NewRecorder()
	router.ServeHTTP(result, request)
	if result.Code != http.StatusOK {
		t.Fatalf("upload status = %d: %s", result.Code, result.Body.String())
	}
	var payload struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(result.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Data.URL == "" {
		t.Fatal("upload response has no URL")
	}
	readResult := httptest.NewRecorder()
	router.ServeHTTP(readResult, httptest.NewRequest(http.MethodGet, payload.Data.URL, nil))
	if readResult.Code != http.StatusOK || !bytes.Equal(readResult.Body.Bytes(), png) {
		t.Fatalf("read status = %d; body = %x", readResult.Code, readResult.Body.Bytes())
	}
}
