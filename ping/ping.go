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
package ping

import (
	"crypto/rand"
	"log/slog"
	"math"
	"net"

	"github.com/LibSEA/mixnet/session"
	"github.com/flynn/noise"
)

type Options struct {
}

func Run(opts Options) int {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		slog.Error("error connecting", "error", err)
		return 1
	}

	cs := noise.NewCipherSuite(noise.DH25519, noise.CipherChaChaPoly, noise.HashBLAKE2b)

	kp, err := cs.GenerateKeypair(rand.Reader)
	if err != nil {
		slog.Error("error generating keypair. panicking.", "error", err)
		return 1
	}

	s := session.New(conn, cs, kp)

	var buf = make([]byte, math.MaxInt16)
	defer s.Close()

	err = s.ClientHandshake(buf)
	if err != nil {
		slog.Error("failed handshake", "error", err)
		return 1
	}

	err = s.WriteMessage(buf, []byte("ping"))
	if err != nil {
		slog.Error("failed write", "error", err)
		return 1
	}

	return 0
}
