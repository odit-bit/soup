package soup

// here lies convinient function for server or client needed of soupbinTCP protocol flow ,
// included login, authentication , heartbeat mechanism.
// user provide it's own implementation of handler func for their specific requirement

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"time"
)

// auth if username and password its matched
type authUserFunc func(username, password string) bool

// auth if session or sequence is avaialable
type authSessionFunc func(rSess, rSeq string) (session string, sequence string, ok bool)

// auth is soupbinTCP specific flow for authenticate process.
// it will read login-request packet and start authentication by given func
// and send response back login-accept or login-rejected packet
func DefaulAuth(rw io.ReadWriter, au authUserFunc, as authSessionFunc) error {
	lr := LoginRequest{}
	if _, err := lr.ReadFrom(rw); err != nil {
		return err
	}

	un := string(bytes.TrimSpace(lr.Username[:]))
	pass := string(bytes.TrimSpace(lr.Password[:]))
	if ok := au(un, pass); !ok {
		lr := LoginReject{
			Reason: 'A',
		}
		if _, err := lr.WriteTo(rw); err != nil {
			return err
		}
		return errors.New("client authenticate failed")
	}

	rSess := string(bytes.TrimSpace(lr.Session[:]))
	rSeq := string(bytes.TrimSpace(lr.Sequence[:]))
	session, sequence, ok := as(rSess, rSeq)
	if !ok {
		lr := LoginReject{
			Reason: 'S',
		}
		if _, err := lr.WriteTo(rw); err != nil {
			return err
		}
		return errors.New("requested session or sequence not available")
	}

	la := LoginAccept{
		Session:  [10]byte([]byte(fmt.Sprintf("%10s", session))),
		Sequence: [20]byte([]byte(fmt.Sprintf("%20s", sequence))),
	}
	_, err := la.WriteTo(rw)
	return err
}

// advancing the read deadline of other node of connection
// while (writer) Conn is idle by sending heartbeat packet
func DefaultHeartbeat(w io.Writer, reset <-chan time.Duration, hbType uint8) {
	defer func() {
		log.Println("exit pinger")
	}()

	var interval time.Duration = 1 * time.Second

	hb := []byte{0, 1, hbType}
	timer := time.NewTimer(interval)
	log.Println("start hb", interval)
	for {
		select {
		case <-timer.C:
			_, err := w.Write(hb)
			if err != nil {
				// log.Println("[DEBUG] pinger:", err)
				return
			}
		case newInterval, ok := <-reset:
			if !timer.Stop() {
				<-timer.C
			}
			if !ok {
				return
			}
			interval = newInterval
		}
		_ = timer.Reset(interval)
	}
}

// // launch go routine in the background to read message from conn and
// // respons according to message type concurrently until it failed
// func event(conn net.Conn, readErr chan<- error, reset chan<- time.Duration) {

// 	defer func() {
// 		log.Println("exit Read-Loop")
// 		close(readErr)
// 	}()

// 	/*read block*/
// 	reader := bufio.NewReader(conn)
// 	for {
// 		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
// 		b, err := reader.Peek(3)
// 		if err != nil {
// 			log.Println("DEBUG PEEK", err)
// 			readErr <- err
// 			break
// 		}

// 		switch b[2] {
// 		case 'L': //auth
// 			if err := auth(reader, conn); err != nil {
// 				readErr <- err
// 				return
// 			}
// 		case 'R': //client heartbeat
// 			reader.Discard(3)
// 			fmt.Println("hb", b)
// 		default:
// 			log.Println("DEBUG unknown packet", b[2])
// 		}
// 		// reset <- 1 * time.Second
// 	}
// }

// // read incoming packet from client
// func readPacket(reader *bufio.Reader) error {
// 	/*read block*/
// 	// reader := bufio.NewReader(conn)
// 	// conn.SetReadDeadline(time.Now().Add(2 * time.Second))

// 	b, err := reader.Peek(3)
// 	if err != nil {
// 		log.Println("DEBUG PEEK", err)
// 		return err
// 	}

// 	switch b[2] {
// 	case 'R': //client heartbeat
// 		reader.Discard(3)
// 		fmt.Println("hb", b)
// 	case 'O':
// 		reader.Discard(3)
// 		fmt.Println("[DEBUG] receive logout")
// 		return nil
// 	default:
// 		log.Println("DEBUG unknown packet", b[2])
// 		return errors.New("unknown packet")
// 	}
// 	return nil
// }
