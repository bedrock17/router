package router

import "net/http"

// Context : URL 패턴 처리를 위한 구조체
type Context struct {
	Param map[string]interface{}

	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

type HandlerFunc func(*Context)
