package main

import (
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"fmt"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

var blackHole = make([]byte, 8192*1024) // 8 MiB

const maxDataSize = 20 * 1024 * 1024 * 1024 // 20G

func main() {

	r := gin.New()
	// catch all panics and log them
	// r.GET("/download", handler)
	r.Any("/*any", handler)
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handler(c *gin.Context) {
	// if c.Request.URL.Path == "/" {
	c.Header("content-type", "text/html; charset=UTF-8")
	// }
	acceptEncodingsStr := c.GetHeader("accept-encoding")
	acceptEncodings := strings.Split(acceptEncodingsStr, ", ")
	if len(acceptEncodings) == 0 {
		return
	}
	switch acceptEncodings[0] {
	case "gzip":
		fmt.Println("gzip")
		c.Header("Content-Encoding", "gzip")
		writer, err := gzip.NewWriterLevel(c.Writer, gzip.BestSpeed)
		if err != nil {
			fmt.Println(err)
			return
		}
		totalWritten := 0
		for totalWritten < maxDataSize {
			n, err := writer.Write(blackHole)
			if err != nil {
				fmt.Println(err)
				return
			}
			totalWritten += n
		}
		writer.Close()
		return
	case "deflate":
		fmt.Println("deflate")
		c.Header("Content-Encoding", "deflate")
		writer, err := flate.NewWriter(c.Writer, flate.BestSpeed)
		if err != nil {
			fmt.Println(err)
			return
		}
		totalWritten := 0
		for totalWritten < maxDataSize {
			n, err := writer.Write(blackHole)
			if err != nil {
				fmt.Println(err)
				return
			}
			totalWritten += n
		}
		writer.Close()
		return
	case "br":
		fmt.Println("br")
		c.Header("Content-Encoding", "br")
		writer := brotli.NewWriterLevel(c.Writer, brotli.BestSpeed)
		totalWritten := 0
		for totalWritten < maxDataSize {
			n, err := writer.Write(blackHole)
			if err != nil {
				fmt.Println(err)
				return
			}
			totalWritten += n
		}
		writer.Close()
	case "compress":
		fmt.Println("compress")
		c.Header("Content-Encoding", "compress")
		writer := lzw.NewWriter(c.Writer, lzw.LSB, 8)
		totalWritten := 0
		for totalWritten < maxDataSize {
			n, err := writer.Write(blackHole)
			if err != nil {
				fmt.Println(err)
				return
			}
			totalWritten += n
		}
		writer.Close()
	default:
		return
	}
}
