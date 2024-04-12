package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
)

var defaultAddress = ":5001"

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	ln        net.Listener
	peers     map[*Peer]bool
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan []byte
}

func NewServer(config Config) *Server {
	if len(config.ListenAddress) == 0 {
		config.ListenAddress = defaultAddress
	}
	return &Server{
		Config:    config,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan []byte),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", defaultAddress)
	if err != nil {
		return err
	}
	s.ln = ln

	// Start listening to the channels before accepting any requests
	go s.channelLoop()

	slog.Info("server running", "listenAddr", s.ListenAddress)
	return s.acceptLoop()
}

// channelLoop listens to all the messages on the Server channels eg: AddPeer
func (s *Server) channelLoop() {
	for {
		select {
		// Block in case a new peer is added
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		case rawMsg := <-s.msgCh:
			if err := s.handleRawMessage(rawMsg); err != nil {
				slog.Error("raw message error", "err", err)
			}
		}
	}
}

func (s *Server) handleRawMessage(rawMsg []byte) error {
	fmt.Print(string(rawMsg))
	return nil
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept failed", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("new peer added", "remoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	config := Config{}
	server := NewServer(config)
	log.Fatal(server.Start())
}
