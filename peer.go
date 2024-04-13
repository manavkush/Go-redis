package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgCh chan<- Message
}

func NewPeer(conn net.Conn, msgCh chan<- Message) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

func (peer *Peer) readLoop() error {
	rd := resp.NewReader(peer.conn)
	for {
		// Reads a single command
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSET:
					if len(v.Array()) != 3 {
						return fmt.Errorf("Invalid number of arguments in the SET command. ", "err", err)
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					peer.msgCh <- Message{
						cmd:  cmd,
						peer: peer,
					}
				case CommandGET:
					if len(v.Array()) != 2 {
						return fmt.Errorf("Invalid number of arguments in the GET command. ", "err", err)
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}
					peer.msgCh <- Message{
						cmd:  cmd,
						peer: peer,
					}
				}
			}
		}
	}
	return nil
}

func (peer *Peer) Send(msgData []byte) (int, error) {
	return peer.conn.Write(msgData)
}
