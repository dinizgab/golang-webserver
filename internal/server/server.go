package server

import (
	"fmt"
	"log"
	"net"
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
    fmt.Println("Hello from handle!")
    // TODO - handle the connection based on the protocol
}
