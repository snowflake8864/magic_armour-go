package authorization

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w ResponseBody) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w ResponseBody) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
