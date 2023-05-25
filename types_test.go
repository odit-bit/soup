package soup

import (
	"bytes"
	"testing"
)

// func newloginPacket(username, password, session, sequence string) []byte {
// 	packet := []byte{}

// 	header := uint16(47)
// 	packet = binary.BigEndian.AppendUint16(packet, header)

// 	typ := []byte{LoginType}
// 	packet = append(packet, typ...)

// 	un := []byte(fmt.Sprintf("%6s", username))
// 	packet = append(packet, un...)

// 	pass := []byte(fmt.Sprintf("%10s", password))
// 	packet = append(packet, pass...)

// 	sess := []byte(fmt.Sprintf("%10s", session))
// 	packet = append(packet, sess...)

// 	seq := []byte(fmt.Sprintf("%20s", sequence))
// 	packet = append(packet, seq...)

// 	return packet
// }

// func Test_newLoginPacket_Func(t *testing.T) {
// 	user := "admin"
// 	password := "12345"
// 	session := "22a10"
// 	seq := "1"

// 	packet := newloginPacket(user, password, session, seq)

// 	if len(packet) != 47 {
// 		t.Errorf("got %v, expect %v", len(packet), 47)
// 	}

// 	if !bytes.Equal(packet[1:7], []byte(" admin")) {
// 		t.Errorf("got %v, expect %v", string(packet[1:7]), " admin")
// 	}

// 	if !bytes.Equal(packet[7:17], []byte("     12345")) {
// 		t.Errorf("got %v, expect %v", string(packet[7:17]), "     12345")
// 	}

// 	if !bytes.Equal(packet[17:27], []byte("     22a10")) {
// 		t.Errorf("got %v, expect %v", string(packet[17:27]), "     22a10")
// 	}

// 	if !bytes.Equal(packet[27:47], []byte("                   1")) {
// 		t.Errorf("got %v, expect %v", string(packet[27:47]), "                   1")
// 	}
// }

// func Test_loginData_ReadFrom(t *testing.T) {
// 	//setup
// 	user := "admin"
// 	password := "12345"
// 	session := "22a10"
// 	seq := "1"
// 	packet := newloginPacket(user, password, session, seq)
// 	network := bytes.NewReader(packet)

// 	// test ReadFrom
// 	reader := bufio.NewReader(network)
// 	_, _ = reader.Discard(3)

// 	var ld LoginData
// 	n, err := ld.ReadFrom(reader)

// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if n != int64(46) {
// 		t.Errorf("got %v, expected %v", n, 47)
// 	}

// 	// test content
// 	content := ld.Content()
// 	if content[0] != user && content[1] != password && content[2] != session && content[3] != seq {
// 		t.Errorf("got %v", content)
// 	}
// }

// func Test_loginData_WriteTo(t *testing.T) {
// 	//setup
// 	user := "admin"
// 	password := "12345"
// 	session := "22a10"
// 	seq := "1"
// 	packet := newloginPacket(user, password, session, seq)

// 	ld := NewLoginData(user, password, session, seq)

// 	//TEST WriteTo
// 	var network bytes.Buffer
// 	n, err := ld.WriteTo(&network)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if n != 49 {
// 		t.Errorf("got %v ,expect %v", n, 49)
// 	}
// 	actual := network.Bytes()
// 	if !bytes.Equal(actual, packet) {
// 		t.Errorf("bytes not equalgot %v, expect %v", actual, packet)
// 	}
// }

func Test_Sequence(t *testing.T) {
	var dest bytes.Buffer

	msg := []byte("message")
	expect := []byte{0, 8, 'S', 'm', 'e', 's', 's', 'a', 'g', 'e'}
	p := SequencePacket{
		Value: msg,
	}

	n, err := p.WriteTo(&dest)
	if err != nil {
		t.Fatal(err)
	}
	if n != 10 {
		t.Fatal(n)
	}

	if string(expect) != dest.String() {
		t.Fatal(dest.Bytes())
	}

	p2 := &SequencePacket{}

	n2, err := p2.ReadFrom(&dest)
	if err != nil {
		t.Fatal(err)
	}
	if n != 10 {
		t.Fatal(n2)
	}

	if string(p.Value) != string(p2.Value) {
		t.Fatal(string(p2.Value))
	}
}

func Test_LoginRequest(t *testing.T) {
	var dst bytes.Buffer
	lr := newLoginRequest("admin", "12345", "22a10", "1")

	n, err := lr.WriteTo(&dst)
	if err != nil {
		t.Fatal(err)
	}
	if n != 49 {
		t.Fatal(n)
	}

	actual := dst.Bytes()
	if !bytes.Equal(actual[0:2], []byte{0, 47}) {
		t.Fatal(actual[0:2])
	}
	if !bytes.Equal(actual[2:3], []byte{'L'}) {
		t.Fatal(actual[2:3])
	}
	if !bytes.Equal(bytes.TrimSpace(actual[3:9]), []byte("admin")) {
		t.Fatal(actual[2:3])
	}
	if !bytes.Equal(bytes.TrimSpace(actual[9:19]), []byte("12345")) {
		t.Fatal(actual[9:19])
	}
	if !bytes.Equal(bytes.TrimSpace(actual[19:29]), []byte("22a10")) {
		t.Fatal(actual[19:29])
	}
	if !bytes.Equal(bytes.TrimSpace(actual[29:49]), []byte("1")) {
		t.Fatal(actual[29:49])
	}

	lr2 := &LoginRequest{}
	_, err = lr2.ReadFrom(&dst)
	if err != nil {
		t.Fatal(err)
	}
	switch {
	case !bytes.Equal(lr.Username[:], lr2.Username[:]):
		t.Fatal("username")
	case !bytes.Equal(lr.Password[:], lr2.Password[:]):
		t.Fatal("password")
	case !bytes.Equal(lr.Session[:], lr2.Session[:]):
		t.Fatal("session")
	case !bytes.Equal(lr.Sequence[:], lr2.Sequence[:]):
		t.Fatal("sequence")

	}
}

func Benchmark_Sequence(b *testing.B) {
	packet := bytes.Repeat([]byte{'X'}, 65537)
	packet[0], packet[1], packet[2] = uint8(255), uint8(255), uint8('S')
	msg := packet

	dst := bytes.NewReader(msg)
	p := new(SequencePacket)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		n, err := p.ReadFrom(dst)
		if err != nil {
			b.Fatal(err, n)
		}
		if n != 65537 {
			b.Fatal(n)
		}

		dst.Reset(msg)
	}
}
