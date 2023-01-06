package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

func ContextBody() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 针对/api/v1/webhook接口设置body文本存储供WebhookService服务打印使用
		fmt.Println(context.Request.URL.Path)
		if context.Request.URL.Path == "/api/v1/webhook" {
			byts, err := PeekRequest(context.Request)
			if err != nil {
				log.Printf("Middleware ContextBody PeakRequest failed, err: %s\n", err)
			} else {
				context.Set("body", byts)
			}
		}
		context.Next()
	}
}

func PeekRequest(request *http.Request) ([]byte, error) {
	if request.Body != nil {
		byts, err := io.ReadAll(request.Body) // io.ReadAll as Go 1.16, below please use ioutil.ReadAll
		if err != nil {
			return nil, err
		}
		request.Body = io.NopCloser(bytes.NewReader(byts))
		return byts, nil
	}
	return make([]byte, 0), nil
}
