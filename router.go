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
	DevMode  bool
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

	//개발모드라면 디버깅에 용이한 정보를 출력하도록 핸들러를 감싼다

	if router.DevMode {
		h = RecoverHandler(LogHandler(h))
	}

	//http 메서드로 등록된 맵에 URL 패턴과 핸들러 함수 등록
	m[pattern] = h
}

func match(pattern, path string) (bool, map[string]string) {
	if pattern == path {
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
	for pattern, handler := range router.Handlers[req.Method] {
		if ok, params := match(pattern, req.URL.Path); ok {
			//요청 URL 에 해당하는 핸들러 수행
			// debug code
			// fmt.Println(pattern, params)

			c := Context{
				Param:          make(map[string]interface{}),
				ResponseWriter: res,
				Request:        req,
			}

			for k, v := range params {
				c.Param[k] = v
			}

			if router.DevMode {
				//개발모드 일 땐 Access-Control-Allow-Origin 을 모두 허용
				c.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
			}

			handler(&c)
			return
		}
	}

	http.NotFound(res, req)
	return
}
