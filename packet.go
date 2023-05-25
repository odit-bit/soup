package soup

import (
	"bytes"
	"encoding/binary"
	"io"
)

// extend read and write functionality of underlying reader and writer
// so that can read and write soupbin logical packet
type Packet struct {
	w   io.Writer
	r   io.Reader
	buf *bytes.Buffer
}

func (pw *Packet) Write(p []byte) (int, error) {
	var n int
	//write header (payload length 2 byte size)
	size := uint16(len(p)) // 2 byte
	if err := binary.Write(pw.w, binary.BigEndian, size); err != nil {
		return n, err
	}
	n += 2

	// write the payload
	o, err := pw.w.Write(p)
	if err != nil {
		n += o
		return n, err
	}
	n += o
	return n, nil

}

// read the next packet and copy to internal buffer
func (pr *Packet) next() (int64, error) {
	var n int64
	size, err := readSizePacket(pr.r)
	if err != nil {
		return n, err
	}

	if pr.buf.Cap() < int(size) {
		pr.buf.Grow(int(size))
	}
	n += 2

	o, err := io.CopyN(pr.buf, pr.r, int64(size))
	if err != nil {
		n += o
		return n, err
	}
	n += o
	return n, nil
}

// implement reader interface
func (pr *Packet) Read(p []byte) (int, error) {
	_, err := pr.next()
	if err != nil {
		return 0, err
	}

	return pr.buf.Read(p)
}

// return the next packet
func (pr *Packet) NextPacket() ([]byte, error) {
	size, err := readSizePacket(pr.r)
	if err != nil {
		return nil, err
	}

	payload := make([]byte, size)
	_, err = pr.r.Read(payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// read the header of packet that consisted the size of payload
func readSizePacket(r io.Reader) (int, error) {
	var size uint16 // 2 byte size
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return 0, err
	}

	return int(size), nil
}

type PacketIO struct{}
