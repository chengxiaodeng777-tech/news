package system

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"base/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

var articleImageName = regexp.MustCompile(`^[a-f0-9]{32}\.(png|jpg|gif|webp)$`)

func articleUploadDirectory() string {
	if directory := os.Getenv("ARTICLE_UPLOAD_DIR"); directory != "" {
		return directory
	}
	return filepath.Join("uploads", "articles")
}

func UploadArticleImage(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择图片文件")
		return
	}
	input, err := header.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无法读取图片")
		return
	}
	defer input.Close()

	var sample [512]byte
	count, err := input.Read(sample[:])
	if err != nil && err != io.EOF {
		response.Error(c, http.StatusBadRequest, "无法读取图片")
		return
	}
	if _, err = input.Seek(0, io.SeekStart); err != nil {
		response.Error(c, http.StatusBadRequest, "无法读取图片")
		return
	}

	extensions := map[string]string{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
		"image/gif":  ".gif",
		"image/webp": ".webp",
	}
	extension, ok := extensions[http.DetectContentType(sample[:count])]
	if !ok {
		response.Error(c, http.StatusBadRequest, "仅支持 PNG、JPEG、GIF 和 WebP 图片")
		return
	}

	directory := articleUploadDirectory()
	if err := os.MkdirAll(directory, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, "无法创建图片目录")
		return
	}
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		response.Error(c, http.StatusInternalServerError, "无法生成图片名称")
		return
	}
	name := hex.EncodeToString(random[:]) + extension
	output, err := os.OpenFile(filepath.Join(directory, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "无法保存图片")
		return
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(filepath.Join(directory, name))
		response.Error(c, http.StatusInternalServerError, "无法保存图片")
		return
	}

	response.Success(c, gin.H{"url": "/api/uploads/articles/" + name})
}

func ServeArticleImage(c *gin.Context) {
	name := c.Param("name")
	if !articleImageName.MatchString(name) {
		c.Status(http.StatusNotFound)
		return
	}
	path := filepath.Join(articleUploadDirectory(), name)
	if _, err := os.Stat(path); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("X-Content-Type-Options", "nosniff")
	c.File(path)
}
