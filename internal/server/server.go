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
}

func New(host string, port int) *Server {
	return &Server{
		Host: host,
		Port: port,
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

	fmt.Println(request)
	// TODO - Create handle functions based on req method and path
	res := response.New(200, map[string]string{"test": "test"}, "Hello World")

	conn.Write([]byte(res.String()))
}
