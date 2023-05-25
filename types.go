package soup

/*
	protocol
	1.client dial with the auth packet
	2.server received auth packet and response back with login-accept or login-reject packet.
	3.client receiver response , ready to receive message
	4.server sending sequence message
*/

import (
	"encoding/binary"
	"fmt"
	"io"
)

const ()
const (
	LoginType             = 'L'
	LoginAcceptType       = 'A'
	LoginRejectType uint8 = 'R'

	ServerHbType = 'H'
)

type RejectReason uint8

type LoginAccept struct {
	Session  [10]byte // requested session ID/number
	Sequence [20]byte // expected sequence number of sequence message
}

func NewLoginAccept(session, sequence string) *LoginAccept {
	la := &LoginAccept{
		Session:  [10]byte([]byte(fmt.Sprintf("%10s", session))),
		Sequence: [20]byte([]byte(fmt.Sprintf("%20s", sequence))),
	}

	return la
}

func (la *LoginAccept) ReadFrom(r io.Reader) (int64, error) {
	var size uint16
	var n int64
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return n, err
	}

	var typ uint8
	if err := binary.Read(r, binary.BigEndian, &typ); err != nil {
		return 2, err
	}

	if err := binary.Read(r, binary.BigEndian, la); err != nil {
		return 3, err
	}

	n = int64(size) + int64(2)
	return n, nil
}

func (la *LoginAccept) WriteTo(w io.Writer) (int64, error) {
	size := uint16(21)
	if err := binary.Write(w, binary.BigEndian, size); err != nil {
		return 0, err
	}

	if err := binary.Write(w, binary.BigEndian, uint8('A')); err != nil {
		return 2, err
	}

	if err := binary.Write(w, binary.BigEndian, la); err != nil {
		return 3, err
	}
	return int64(21), nil
}

type LoginReject struct {
	Reason RejectReason
}

func (lr *LoginReject) WriteTo(w io.Writer) (int64, error) {
	p := []byte{0, 2, 'J', byte(lr.Reason)}
	n, err := w.Write(p)
	if err != nil {
		return int64(n), err
	}

	return int64(n), nil
}

func (lr *LoginReject) ReadFrom(r io.Reader) (int64, error) {
	p := make([]byte, 4)
	_, err := io.ReadFull(r, p)
	if err != nil {
		return 0, err
	}
	lr.Reason = RejectReason(p[3])
	return 4, nil
}

type LoginRequest struct {
	// Length   uint16
	// Typ      uint8
	Username [6]byte
	Password [10]byte
	Session  [10]byte
	Sequence [20]byte
}

func newLoginRequest(username, password, session, sequence string) *LoginRequest {
	lr := LoginRequest{
		// Length:   46 + 1,
		// Typ:      'L',
		Username: [6]byte([]byte(fmt.Sprintf("%6s", username))),
		Password: [10]byte([]byte(fmt.Sprintf("%10s", password))),
		Session:  [10]byte([]byte(fmt.Sprintf("%10s", session))),
		Sequence: [20]byte([]byte(fmt.Sprintf("%20s", sequence))),
	}
	return &lr
}

func (lr *LoginRequest) WriteTo(w io.Writer) (int64, error) {
	size := uint16(47)
	if err := binary.Write(w, binary.BigEndian, size); err != nil {
		return 0, err
	}

	if err := binary.Write(w, binary.BigEndian, uint8('L')); err != nil {
		return 2, err
	}

	if err := binary.Write(w, binary.BigEndian, lr); err != nil {
		return 3, err
	}
	return int64(49), nil
}

func (lr *LoginRequest) ReadFrom(r io.Reader) (int64, error) {
	var size uint16
	var n int64
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return n, err
	}

	var typ uint8
	if err := binary.Read(r, binary.BigEndian, &typ); err != nil {
		return 2, err
	}

	if err := binary.Read(r, binary.BigEndian, lr); err != nil {
		return 3, err
	}

	n = int64(size) + int64(2)
	return n, nil
}

type SequencePacket struct {
	Value []byte
}

func (sp *SequencePacket) WriteTo(w io.Writer) (int64, error) {

	size := uint16(1 + len(sp.Value))
	if err := binary.Write(w, binary.BigEndian, size); err != nil {
		return 0, err
	}

	if err := binary.Write(w, binary.BigEndian, uint8('S')); err != nil {
		return 0, err
	}
	if n, err := w.Write(sp.Value); err != nil {
		return int64(n), err
	}

	return int64(2 + size), nil
}

func (sp *SequencePacket) ReadFrom(r io.Reader) (int64, error) {
	var size uint16
	var n int64
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return n, err
	}
	n += 2

	payload := make([]byte, size)
	o, err := io.ReadFull(r, payload)
	if err != nil {
		return n, err
	}

	n += int64(o)
	sp.Value = payload[1:]
	return n, nil
}

// ///
