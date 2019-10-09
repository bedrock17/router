package router

import (
	"net/http"
	"strings"
)

//Handler : http요청을 처리하기 위한 인터페이스
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

//Router : http 요청을 처리하는 핸드러를 담고있는 map
//[method][path] 구조이다.
type Router struct {
	Handlers map[string]map[string]HandlerFunc
}

//HandleFunc : Router에 http 핸들러를 추가한다.
func (router *Router) HandleFunc(method, pattern string, h HandlerFunc) {
	//http 메서드로 등록된 맵이 있는지 확인
	m, ok := router.Handlers[method]
	if !ok {
		//등록된 맵이 없다면 새로운 맵을 생성
		m = make(map[string]HandlerFunc)

		router.Handlers[method] = m
	}

	//http 메서드로 등록된 맵에 URL 패턴과 핸들러 함수 등록
	m[pattern] = h
}

func match(pattern, path string) (bool, map[string]string) {

	// fmt.Println("[debug] pattern :", pattern, pattern[:len(pattern)-1], "!! path", path, strings.HasSuffix(pattern, "*"), strings.HasPrefix(path, pattern[:len(pattern)-1]))

	if pattern == path {
		return true, nil
	} else if strings.HasSuffix(pattern, "*") && strings.HasPrefix(path, pattern[:len(pattern)-1]) { // ~~/* 형태로 끝나는 URI를 모두 처리하기 위해..
		return true, nil
	}
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	if len(patterns) != len(paths) {

		return false, nil
	}

	params := make(map[string]string)

	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
		case len(patterns[i]) > 0 && patterns[i][0] == ':':

			params[patterns[i][1:]] = paths[i]
		default:

			return false, nil
		}
	}
	return true, params
}

func (router *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for pattern, handle := range router.Handlers[req.Method] {
		if ok, params := match(pattern, req.URL.Path); ok {
			//요청 URL 에 해당하는 핸들러 수행
			// debug code
			// fmt.Println(pattern, params)

			c := Context{
				Params:         make(map[string]interface{}),
				ResponseWriter: res,
				Request:        req,
			}

			for k, v := range params {
				c.Params[k] = v
			}

			handle(&c)
			return
		}
	}

	http.NotFound(res, req)
	return
}

func (router *Router) handler() HandlerFunc {
	return func(c *Context) {
		// http 메서드에 맞는 모든 handers를 반복하며 요청 URL에 해당하는 handler를 찾음
		// Todo Handler
		for pattern, handlerFunc := range router.Handlers[c.Request.Method] {

			if ok, params := match(pattern, c.Request.URL.Path); ok {
				for k, v := range params {
					c.Params[k] = v
				}
				// 요청 URL에 해당하는 handler 수행
				handlerFunc(c)
				return
			}
		}

		// 요청 URL에 해당하는 handler를 찾지 못하면 NotFound 에러 처리
		http.NotFound(c.ResponseWriter, c.Request)
		return
	}
}
