package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"path/filepath"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置CORS相关的响应头
		// 这里允许所有来源的请求（出于安全考虑，你应该只允许特定的源）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 检查请求的方法是否为OPTIONS，因为浏览器在发送实际请求之前会先发送一个预检请求（preflight request）
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}

		// 继续处理后续的处理函数
		c.Next()
	}
}

const (
	staticFilePath = "/usr/local/project/static"
)

func main() {
	// 创建一个默认的Gin引擎
	r := gin.Default()

	r.Use(CORS())
	// 定义文件上传的路由
	r.POST("/file/avatar", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		fromValue := form.Value["from"]
		if fromValue == nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		from := fromValue[0]
		files := form.File["file"]
		// 遍历文件切片
		for _, file := range files {
			f, err := file.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{})
				return
			}
			hash := md5.New()
			_, err = io.Copy(hash, f)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{})
				return
			}
			fileHash := hex.EncodeToString(hash.Sum(nil))
			fileExt := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("%s/%s%s", from, fileHash, fileExt)
			saveFilePath := fmt.Sprintf("%s/%s", staticFilePath, newFileName)
			if err := c.SaveUploadedFile(file, saveFilePath); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"fileHash": newFileName,
				})
			}
		}
	})

	// 启动服务并监听在 8080 端口
	r.Run(":8889")
}
