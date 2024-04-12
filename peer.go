package main

import (
	"log/slog"
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan<- []byte
}

func NewPeer(conn net.Conn, msgCh chan<- []byte) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

func (peer *Peer) readLoop() error {
	buf := make([]byte, 1024)
	for {
		n, err := peer.conn.Read(buf)
		if err != nil {
			slog.Error("peer read error.", "err", err)
			return err
		}

		msgBuf := make([]byte, n)
		copy(msgBuf, buf)
		peer.msgCh <- msgBuf
	}
}
