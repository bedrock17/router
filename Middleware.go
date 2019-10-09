package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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

//StaticHandler 정적파일 처리함수
func StaticHandler(next HandlerFunc) HandlerFunc {
	var (
		dir = http.Dir(".")
		// indexFile = "index.html"
	)
	return func(c *Context) {

		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			next(c)
			return
		}

		file := c.Request.URL.Path

		f, err := dir.Open(file)
		if err != nil {
			next(c)
			return
		}

		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			next(c)
			return
		}

		if fi.IsDir() {

			if !strings.HasSuffix(c.Request.URL.Path, "/") {
				http.Redirect(c.ResponseWriter, c.Request, c.Request.URL.Path+"/", http.StatusFound)
				return
			}
		}

		// file = path.Join(file, indexFile)

		f, err = dir.Open(file)
		if err != nil {
			fmt.Println(err)
			next(c)
			return
		}

		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() {
			next(c)
			return
		}

		http.ServeContent(c.ResponseWriter, c.Request, file, fi.ModTime(), f)
	}

}
