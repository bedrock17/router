package router

import (
	"log"
	"net/http"
	"time"
)

//Middleware : 미들웨어 함수타입
type Middleware func(next HandlerFunc) HandlerFunc

//LogHandler : Log를 남겨주는 핸들러
func LogHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		next(c)
		t := time.Now()
		log.Printf("[%s] %q %v\n",
			c.Request.Method,
			c.Request.URL.String(),
			time.Now().Sub(t))

	}
}

//RecoverHandler : 내부 오류를 보여준다
func RecoverHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(c.ResponseWriter,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next(c)
	}
}
