package main

import (
	"flag"
	"log"
	"log/slog"
	"net"
)

var defaultListAddr = ":5001"

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	ln        net.Listener
	peers     map[*Peer]bool
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan Message

	kv *KV
}

func NewServer(config Config) *Server {
	if len(config.ListenAddress) == 0 {
		config.ListenAddress = defaultListAddr
	}
	return &Server{
		Config:    config,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", defaultListAddr)
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
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw message error", "err", err)
			}
		}
	}
}

func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case SetCommand:
		s.kv.Set(v.key, v.val)
		// slog.Info("Someone wants to set a key inside the redis", "key", v.key, "value", v.val)
	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			slog.Error("Error in getting the value for key.", "key", v.key)
			return nil
		}
		_, err := msg.peer.Send(val)
		if err != nil {
			slog.Error("Error in sedning message to peer", "err", err)
		}
	}
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
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
	listenAddr := flag.String("listenAddr", defaultListAddr, "address for the server to listen on")
	flag.Parse()
	config := Config{
		ListenAddress: *listenAddr,
	}
	server := NewServer(config)
	log.Fatal(server.Start())
}
