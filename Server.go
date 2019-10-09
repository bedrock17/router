package router

import (
	"net/http"
)

//Server : 미들웨어를 유연하게 처리하기 위한 구조체
type Server struct {
	*Router
	middlewares  []Middleware
	startHandler HandlerFunc
}

//NewServer : 새로운 서버를 생성한다.
func NewServer() *Server {
	r := &Router{Handlers: make(map[string]map[string]HandlerFunc), DevMode: false} //필요하면 DevMode는 따로 활성화 한다.
	s := &Server{Router: r}
	return s
}

//AppendMiidleWare : 미들웨어를 추가한다.
func (s *Server) AppendMiidleWare(m Middleware) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Context 생성
	c := &Context{
		Params:         make(map[string]interface{}),
		ResponseWriter: w,
		Request:        r,
	}
	for k, v := range r.URL.Query() {
		c.Params[k] = v[0]
	}
	s.startHandler(c)
}

//Run : 웹서버를 실행한다.
func (s *Server) Run(addr string) {
	// startHandler를 라우터 핸들러 함수로 지정
	s.startHandler = s.Router.handler()

	// 등록된 미들웨어를 라우터 핸들러 앞에 하나씩 추가
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.startHandler = s.middlewares[i](s.startHandler)
	}

	// 웹 서버 시작
	if err := http.ListenAndServe(addr, s); err != nil {
		panic(err)
	}
}
