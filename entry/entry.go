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
package entry

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"math"
	"net"
	"os"
	"strconv"

	"github.com/LibSEA/mixnet/session"
	"github.com/flynn/noise"
)

type Options struct {
	Port uint16
	Host string
}

type cmd struct {
	logger *slog.Logger
}

func (c *cmd) handle(s *session.Session) {
	defer func() { _ = s.Close() }()

	var buf = make([]byte, math.MaxInt16)

	err := s.ServerHandshake(buf)
	if err != nil {
		c.logger.Warn("ServerHandshake failed.", "error", err)
		return
	}

	for {
		msg, err := s.ReadMessage(buf)
		if err != nil {
			c.logger.Warn("ReadMessage failed", "error", err)
			return
		}
		fmt.Println(string(msg))
	}
}

func Run(opts Options) int {
	var c = cmd{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	ln, err := net.Listen(
		"tcp",
		net.JoinHostPort(
			opts.Host,
			strconv.Itoa(int(opts.Port))),
	)
	if err != nil {
		c.logger.Error(
			"couldn't listen",
			"port", opts.Port,
			"host", opts.Host,
		)
		return 1
	}

	cs := noise.NewCipherSuite(
		noise.DH25519, noise.CipherChaChaPoly, noise.HashBLAKE2b,
	)
	kp, err := cs.GenerateKeypair(rand.Reader)

	if err != nil {
		c.logger.Error("error generating keypair.", "error", err)
		return 1
	}

	cf := 0

	c.logger.Info("started")

	for {
		if cf > 10 {
			c.logger.Error("failed calling accept to many times.")
			return 1
		}
		conn, err := ln.Accept()
		if err != nil {
			c.logger.Error("error calling accept", "error", err)
			cf++
			continue
		}
		cf = 0
		go c.handle(session.New(conn, cs, kp))
	}
}
