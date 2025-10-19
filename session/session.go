package session

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/flynn/noise"
)

var cs noise.CipherSuite
var kp noise.DHKey

func init() {
	var err error
	cs = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashBLAKE2b)

	kp, err = cs.GenerateKeypair(rand.Reader)

	if err != nil {
		panic("failed to generate keypair")
	}
}

type Session struct {
	b  []byte
	c  net.Conn
	rx *noise.CipherState
	tx *noise.CipherState
}

/*
The handshake we have is one byte for protocol version, followed by 2 bytes
for size, then bytes for message. For every other message it is 2 bytes for
size then the payload.
*/
func New(conn net.Conn) *Session {
	return &Session{
		b: make([]byte, 655535),
		c: conn,
	}
}

func (s *Session) ReadMessage() ([]byte, error) {
	l, err := s.read()

	if err != nil {
		return nil, err
	}

	return s.rx.Decrypt(nil, nil, s.b[0:l])
}

func (s *Session) WriteMessage(in []byte) error {
	msg, err := s.tx.Encrypt(s.b[:2], nil, in)
	if err != nil {
		return err
	}

	return s.write(len(msg))
}

func (s *Session) read() (uint16, error) {
	_, err := io.ReadFull(s.c, s.b[0:2])
	if err != nil {
		return 0, err
	}

	l := binary.BigEndian.Uint16(s.b[0:2])

	_, err = io.ReadFull(s.c, s.b[0:l])
	if err != nil {
		return 0, err
	}

	return l, nil
}

func (s *Session) write(l int) error {

	binary.BigEndian.PutUint16(s.b, uint16(l-2))

	_, err := s.c.Write(s.b[0:l])

	return err
}

func (s *Session) handshakeRead(hs *noise.HandshakeState) error {
	l, err := s.read()
	if err != nil {
		return err
	}

	_, s.rx, s.tx, err = hs.ReadMessage(nil, s.b[0:l])
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) handshakeWrite(hs *noise.HandshakeState) error {
	var msg []byte
	var err error

	msg, s.tx, s.rx, err = hs.WriteMessage(s.b[:2], nil)
	if err != nil {
		return err
	}

	return s.write(len(msg))

}

func (s *Session) ClientHandshake() error {
	hs, err := noise.NewHandshakeState(noise.Config{
		CipherSuite:   cs,
		Pattern:       noise.HandshakeXX,
		Random:        rand.Reader,
		Initiator:     true,
		StaticKeypair: kp,
	})

	if err != nil {
		return err
	}

	// -> v
	s.b[0] = 0x1
	_, err = s.c.Write(s.b[0:1])
	if err != nil {
		return err
	}

	// -> e
	err = s.handshakeWrite(hs)
	if err != nil {
		return err
	}

	// <- e, dhee, s, dhse
	err = s.handshakeRead(hs)
	if err != nil {
		return err
	}

	// -> s, dhse
	err = s.handshakeWrite(hs)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) ServerHandshake() error {
	hs, err := noise.NewHandshakeState(noise.Config{
		CipherSuite:   cs,
		Pattern:       noise.HandshakeXX,
		Random:        rand.Reader,
		Initiator:     false,
		StaticKeypair: kp,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// -> v
	_, err = io.ReadFull(s.c, s.b[0:1])
	if err != nil {
		fmt.Println(err)
		return err
	}

	if s.b[0] != 0x1 {
		return errors.ErrUnsupported
	}

	// -> e
	err = s.handshakeRead(hs)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// <- e, dhee, s, dhse
	err = s.handshakeWrite(hs)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// -> s, dhse
	err = s.handshakeRead(hs)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *Session) Close() {
	s.c.Close()
}
