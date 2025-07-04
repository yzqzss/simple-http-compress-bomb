package main

import (
	_ "embed"
	"fmt"

	"github.com/gin-gonic/gin"
)

//go:embed boom.gz
var boom []byte

func main() {
	r := gin.New()
	// catch all panics and log them
	// r.GET("/download", handler)
	r.Any("/*any", fileHandler)
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func fileHandler(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=UTF-8")
	c.Header("Content-Encoding", "gzip")
	c.Writer.WriteHeader(200)
	_, err := c.Writer.Write(boom)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Writer.Flush()
}
