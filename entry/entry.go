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
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/LibSEA/mixnet/session"
)

type Options struct {
	Port string
}

var log *slog.Logger

func init() {
	log = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func handle(s *session.Session) {
	defer s.Close()

	err := s.ServerHandshake()
	if err != nil {
		log.Info("ServerHandshake failed.", "error", err)
		return
	}

	for {
		msg, err := s.ReadMessage()
		if err != nil {
			fmt.Println("REadMessage failed", err)
			return
		}
		fmt.Println(string(msg))
	}
}

func Run(opts Options) {
	ln, err := net.Listen("tcp", opts.Port)
	if err != nil {
		panic("failed to listen")
	}

	cf := 0

	for {
		if cf > 10 {
			panic("error accepting connections")
		}
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("error calling accept", "error", err)
			cf++
			continue
		}
		cf = 0
		go handle(session.New(conn))
	}

}
