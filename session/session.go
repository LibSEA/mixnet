package session

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/flynn/noise"
)

var cs noise.CipherSuite
var kp noise.DHKey
var logger slog.Logger

func init() {
	var err error
	cs = noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashBLAKE2b)

	kp, err = cs.GenerateKeypair(rand.Reader)

	if err != nil {
		slog.Error("error generating keypair. panicking.", "error", err)
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
		return nil, fmt.Errorf("failed to read message in ReadMessage. %w", err)
	}

	b, err := s.rx.Decrypt(nil, nil, s.b[0:l])

	if err != nil {

		return b, fmt.Errorf("failed to decrypt message. %w", err)
	}

	return b, nil
}

func (s *Session) WriteMessage(in []byte) error {
	msg, err := s.tx.Encrypt(s.b[:2], nil, in)
	if err != nil {
		return fmt.Errorf("failed to encrypt data for writing. %w", err)
	}

	err = s.write(len(msg))
	if err != nil {
		return fmt.Errorf("failed to write message data in WriteMessage. %w", err)
	}

	return nil
}

func (s *Session) read() (uint16, error) {
	_, err := io.ReadFull(s.c, s.b[0:2])
	if err != nil {
		return 0, fmt.Errorf("failed to read length message. %w", err)
	}

	l := binary.BigEndian.Uint16(s.b[0:2])

	_, err = io.ReadFull(s.c, s.b[0:l])
	if err != nil {
		return 0, fmt.Errorf("failed to read message data. %w", err)
	}

	return l, nil
}

func (s *Session) write(l int) error {
	binary.BigEndian.PutUint16(s.b, uint16(l-2))

	_, err := s.c.Write(s.b[0:l])
	if err != nil {
		return fmt.Errorf("error writing message. %w", err)
	}

	return nil
}

func (s *Session) handshakeRead(hs *noise.HandshakeState) error {
	l, err := s.read()
	if err != nil {
		return fmt.Errorf("error reading from socket in handshakeRead. %w", err)
	}

	_, s.rx, s.tx, err = hs.ReadMessage(nil, s.b[0:l])
	if err != nil {
		return fmt.Errorf("error calling ReadMessage in hanshake. %w", err)
	}

	return nil
}

func (s *Session) handshakeWrite(hs *noise.HandshakeState) error {
	var msg []byte
	var err error

	msg, s.tx, s.rx, err = hs.WriteMessage(s.b[:2], nil)
	if err != nil {
		return fmt.Errorf("handshakeWrite failed to WriteMessage. %w", err)
	}

	err = s.write(len(msg))
	if err != nil {
		return fmt.Errorf("failed to write to socket in handshakeWrite. %w", err)
	}

	return nil
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
		return fmt.Errorf("failed to create handshake state. %w", err)
	}

	// -> v
	s.b[0] = 0x1
	_, err = s.c.Write(s.b[0:1])
	if err != nil {
		return fmt.Errorf("ClientHandshake -> v failed. %w", err)
	}

	// -> e
	err = s.handshakeWrite(hs)
	if err != nil {
		return fmt.Errorf("ClientHandshake -> e. %w", err)
	}

	// <- e, dhee, s, dhse
	err = s.handshakeRead(hs)
	if err != nil {
		return fmt.Errorf("ClientHandshake <- e, dhee, s, dhse. %w", err)
	}

	// -> s, dhse
	err = s.handshakeWrite(hs)
	if err != nil {
		return fmt.Errorf("ClientHandshake -> s, dhse. %w", err)
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
		return fmt.Errorf(
			"failed to create handshake state in ServerHandshake. %w",
			err,
		)
	}

	// -> v
	_, err = io.ReadFull(s.c, s.b[0:1])
	if err != nil {
		return fmt.Errorf("ServerHandshake -> v failed. %w", err)
	}

	if s.b[0] != 0x1 {
		return fmt.Errorf("unsupported version {}", s.b[0])
	}

	// -> e
	err = s.handshakeRead(hs)
	if err != nil {
		return fmt.Errorf("ServerHandshake -> e failed. %w", err)
	}

	// <- e, dhee, s, dhse
	err = s.handshakeWrite(hs)
	if err != nil {
		return fmt.Errorf("ServerHandshake <- e, dhee, s, dhse failed. %w", err)
	}

	// -> s, dhse
	err = s.handshakeRead(hs)
	if err != nil {
		return fmt.Errorf("ServerHandshake -> s, dhse failed. %w", err)
	}

	return nil
}

func (s *Session) Close() {
	s.c.Close()
}
