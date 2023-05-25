package soup

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	defaultTimer time.Duration = 1 * time.Second
)

type serverConn struct {
	AuthUser    authUserFunc
	AuthSession authSessionFunc

	conn    net.Conn
	r       *bufio.Reader
	reset   chan time.Duration
	readErr chan error
}

func NewServerConn(c net.Conn) *serverConn {
	s := serverConn{
		conn:    c,
		r:       bufio.NewReader(c),
		reset:   make(chan time.Duration),
		readErr: make(chan error),
	}
	return &s
}

func (s *serverConn) Auth(au authUserFunc, as authSessionFunc) error {
	s.AuthUser = au
	s.AuthSession = as
	if err := s.authentication(); err != nil {
		return err
	}

	go DefaultHeartbeat(s.conn, s.reset, uint8('H'))
	go s.readPacket()
	// go event(s.conn, s.readErr, s.reset)
	return nil
}

func (s *serverConn) authentication() error {
	if s.AuthUser == nil {
		return errors.New("handler for AuthUser cannot be nil")
	}
	if s.AuthSession == nil {
		return errors.New("handler for AuthSession cannot be nil")
	}

	return DefaulAuth(s.conn, s.AuthUser, s.AuthSession)
}

func (s *serverConn) Write(msg SequencePacket) error {
	if _, err := msg.WriteTo(s.conn); err != nil {
		return err
	}
	s.reset <- defaultTimer
	return nil
}

// read incoming packet from client
// and response according to the type
func (s *serverConn) readPacket() {
	defer func() {
		log.Println("exit readPacket")
		close(s.readErr)
	}()

	for {
		b, err := s.r.Peek(3)
		if err != nil {
			log.Println("DEBUG PEEK", err)
			s.readErr <- err
			return
		}

		switch b[2] {
		case 'R': //client heartbeat
			s.r.Discard(3)
			fmt.Println("hb", b)
		case 'O':
			s.r.Discard(3)
			fmt.Println("[DEBUG] receive logout")
			return
		default:
			log.Println("DEBUG unknown packet", b[2])
			s.readErr <- errors.New("unknown packet")
			return
		}
	}
}

func (s *serverConn) Close() error {
	err, ok := <-s.readErr
	if ok {
		<-s.readErr
	}
	return err
}
