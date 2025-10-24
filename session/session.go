/*
mixnet - tool to create and manage LibSEA mixnets
Copyright (C) 2025  Liberatory Sofware Engineering Association

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package session

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/flynn/noise"
)

type Session struct {
	l  []byte
	c  io.ReadWriteCloser
	cs noise.CipherSuite
	kp noise.DHKey
	rx *noise.CipherState
	tx *noise.CipherState
}

/*
The handshake we have is one byte for protocol version, followed by 2 bytes
for size, then bytes for message. For every other message it is 2 bytes for
size then the payload.
*/
func New(conn io.ReadWriteCloser, cs noise.CipherSuite, kp noise.DHKey) *Session {
	ses := Session{
		cs: cs,
		kp: kp,
		c:  conn,
	}

	var b = [2]byte{0, 0}

	ses.l = b[:]

	return &ses
}

func (s *Session) Reinit(conn io.ReadWriteCloser, cs noise.CipherSuite, kp noise.DHKey) {
	s.cs = cs
	s.c = conn
	s.kp = kp
	s.rx = nil
	s.tx = nil
}

func (s *Session) ReadMessage(out []byte) ([]byte, error) {
	msg, err := s.read(out)
	if err != nil {
		return nil, fmt.Errorf("failed to read message in ReadMessage. %w", err)
	}

	b, err := s.rx.Decrypt(nil, nil, msg)
	if err != nil {
		return b, fmt.Errorf("failed to decrypt message. %w", err)
	}

	return b, nil
}

func (s *Session) WriteMessage(out []byte, in []byte) error {
	msg, err := s.tx.Encrypt(out[:0], nil, in)
	if err != nil {
		return fmt.Errorf("failed to encrypt data for writing. %w", err)
	}

	err = s.write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message data in WriteMessage. %w", err)
	}

	return nil
}

func (s *Session) read(out []byte) ([]byte, error) {
	var ls = s.l[:]

	_, err := io.ReadFull(s.c, ls)
	if err != nil {
		return nil, fmt.Errorf("failed to read length message. %w", err)
	}

	l := binary.BigEndian.Uint16(ls)

	var ret = out[:l]

	_, err = io.ReadFull(s.c, ret)
	if err != nil {
		return nil, fmt.Errorf("failed to read message data. %w", err)
	}

	return ret, nil
}

func (s *Session) write(payload []byte) error {
	if len(payload) > math.MaxInt16 {
		return fmt.Errorf("message payload too large.")
	}

	var ls = s.l[:]

	binary.BigEndian.PutUint16(ls, uint16(len(payload)))

	_, err := s.c.Write(ls)
	if err != nil {
		return fmt.Errorf("error writing message len. %w", err)
	}
	_, err = s.c.Write(payload)
	if err != nil {
		return fmt.Errorf("error writing message. %w", err)
	}

	return nil
}

func (s *Session) handshakeRead(out []byte, hs *noise.HandshakeState) error {
	msg, err := s.read(out)
	if err != nil {
		return fmt.Errorf("error reading from socket in handshakeRead. %w", err)
	}

	_, s.rx, s.tx, err = hs.ReadMessage(nil, msg)
	if err != nil {
		return fmt.Errorf("error calling ReadMessage in hanshake. %w", err)
	}

	return nil
}

func (s *Session) handshakeWrite(out []byte, hs *noise.HandshakeState) error {
	var msg []byte
	var err error

	msg, s.tx, s.rx, err = hs.WriteMessage(out[:0], nil)
	if err != nil {
		return fmt.Errorf("handshakeWrite failed to WriteMessage. %w", err)
	}

	err = s.write(msg)
	if err != nil {
		return fmt.Errorf("failed to write to socket in handshakeWrite. %w", err)
	}

	return nil
}

func (s *Session) ClientHandshake(out []byte) error {
	hs, err := noise.NewHandshakeState(noise.Config{
		CipherSuite:   s.cs,
		Pattern:       noise.HandshakeXX,
		Random:        rand.Reader,
		Initiator:     true,
		StaticKeypair: s.kp,
	})

	if err != nil {
		return fmt.Errorf("failed to create handshake state. %w", err)
	}

	// -> v
	var v = []byte{0x1}
	_, err = s.c.Write(v)
	if err != nil {
		return fmt.Errorf("ClientHandshake -> v failed. %w", err)
	}

	// -> e
	err = s.handshakeWrite(out, hs)
	if err != nil {
		return fmt.Errorf("ClientHandshake -> e. %w", err)
	}

	// <- e, dhee, s, dhse
	err = s.handshakeRead(out, hs)
	if err != nil {
		return fmt.Errorf("ClientHandshake <- e, dhee, s, dhse. %w", err)
	}

	// -> s, dhse
	err = s.handshakeWrite(out, hs)
	if err != nil {
		return fmt.Errorf("ClientHandshake -> s, dhse. %w", err)
	}

	return nil
}

func (s *Session) ServerHandshake(out []byte) error {
	hs, err := noise.NewHandshakeState(noise.Config{
		CipherSuite:   s.cs,
		Pattern:       noise.HandshakeXX,
		Random:        rand.Reader,
		Initiator:     false,
		StaticKeypair: s.kp,
	})
	if err != nil {
		return fmt.Errorf(
			"failed to create handshake state in ServerHandshake. %w",
			err,
		)
	}

	// -> v
	var v = []byte{0x0}
	_, err = io.ReadFull(s.c, v)
	if err != nil {
		return fmt.Errorf("ServerHandshake -> v failed. %w", err)
	}

	if v[0] != 0x1 {
		return fmt.Errorf("unsupported version %d", v[0])
	}

	// -> e
	err = s.handshakeRead(out, hs)
	if err != nil {
		return fmt.Errorf("ServerHandshake -> e failed. %w", err)
	}

	// <- e, dhee, s, dhse
	err = s.handshakeWrite(out, hs)
	if err != nil {
		return fmt.Errorf("ServerHandshake <- e, dhee, s, dhse failed. %w", err)
	}

	// -> s, dhse
	err = s.handshakeRead(out, hs)
	if err != nil {
		return fmt.Errorf("ServerHandshake -> s, dhse failed. %w", err)
	}

	return nil
}

func (s *Session) Close() {
	s.c.Close()
}
