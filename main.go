package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

// here we will get the feed if someone wants it
// here people suscribe to thongs
func (s *Server) feed(ws *websocket.Conn) {
	fmt.Println("new inncoimming connection from client for orderbook:", ws.RemoteAddr())
	for {
		feedloading := fmt.Sprintf("order feed data : %d\n", time.Now().UnixNano())
		ws.Write([]byte(feedloading))
		time.Sleep(time.Second * 2)
	}

}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new inncoimming connection from client:", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {

	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[:n]

		s.brodcast(msg)
	}
}

func (s *Server) brodcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("error brodcast :", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderbook", websocket.Handler(server.feed))
	http.ListenAndServe(":3000", nil)
}
