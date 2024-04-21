package server

import (
	"fmt"
	"log"
	"net"

	"github.com/dinizgab/golang-webserver/internal/request"
	"github.com/dinizgab/golang-webserver/internal/response"
)

type Server struct {
	Host string
	Port int
    Handlers map[string]func(*request.Request) *response.Response
}

type serverImpl interface {
    Serve()
    AddHandler(path string, handlerFunc func(*request.Request))
    handle(conn net.Conn)
}

func New(host string, port int) *Server {
	return &Server{
		Host: host,
		Port: port,
        Handlers: map[string]func(*request.Request) *response.Response{},
	}
}

func (s *Server) Serve() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Server listening on %s:%d\n", s.Host, s.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}

		go s.handle(conn)
	}
}

func (s *Server) AddHandler(path string, handlerFunc func(*request.Request) *response.Response) {
    s.Handlers[path] = handlerFunc
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		res := response.New(500, map[string]string{}, "Internal Server Error")
		conn.Write([]byte(res.String()))

		return
	}
	req := string(buffer[:n])
	request, err := request.New(req)
	if err != nil {
		res := response.New(500, map[string]string{}, "Internal Server Error")
		conn.Write([]byte(res.String()))

		return
	}

    handler := s.matchHandler(request.Method, request.Path)//s.Handlers[request.Path]
    if handler == nil {
        res := response.New(404, map[string]string{}, "")
		conn.Write([]byte(res.String()))

        return
    }
    
    res := handler(request)
	conn.Write([]byte(res.String()))
}

func (s *Server) matchHandler(method string, path string) func(*request.Request) *response.Response {
    return nil
}
