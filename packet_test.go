package soup

import (
	"bytes"
	"testing"
)

func Test_ReadPacket(t *testing.T) {
	var destination bytes.Buffer
	packet := []byte{0, 5, 'h', 'e', 'l', 'l', 'o', 0, 3, 'a', 's', 'd'}
	expect := []byte{'h', 'e', 'l', 'l', 'o'}
	destination.Write(packet)

	p := Packet{
		w:   nil,
		r:   &destination,
		buf: &bytes.Buffer{},
	}

	b, err := p.NextPacket()
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expect, b) {
		t.Errorf("got %v , expect %v \n", b, expect)
	}
}

func Test_WritePacket(t *testing.T) {
	var destination bytes.Buffer
	packet := []byte{'h', 'e', 'l', 'l', 'o'}
	expect := []byte{0, 5, 'h', 'e', 'l', 'l', 'o'}

	w := Packet{w: &destination}
	n, err := w.Write(packet)
	if err != nil {
		t.Error(err)
	}
	if n != len(expect) {
		t.Errorf("got %v, expect %v", n, len(expect))
	}
}

func Benchmark_Write(b *testing.B) {
	var dst bytes.Buffer
	payload := bytes.Repeat([]byte{'X'}, 65535) //[]byte{'h', 'e', 'l', 'l', 'o'}
	packet := append([]byte{255, 255}, payload...)

	w := Packet{
		w:   &dst,
		r:   nil,
		buf: &bytes.Buffer{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.Write(packet)
		if err != nil {
			b.Fatal(err)
		}
		dst.Reset()
	}
}

func Benchmark_Read(b *testing.B) {
	payload := bytes.Repeat([]byte{'X'}, 65535) //[]byte{'h', 'e', 'l', 'l', 'o'}
	packet := append([]byte{255, 255}, payload...)
	src := bytes.NewReader(packet)

	p := Packet{
		w:   nil,
		r:   src,
		buf: &bytes.Buffer{},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		payload, err := p.NextPacket()
		if err != nil {
			b.Fatal(err)
		}
		_ = payload
		src.Reset(packet)
	}
}
