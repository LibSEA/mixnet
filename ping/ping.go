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
	"fmt"
	"net"

	"github.com/LibSEA/mixnet/session"
)

type Options struct {
}

func Run(opts Options) {
	fmt.Println("ping called")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic("error connecting")
	}

	s := session.New(conn)

	err = s.ClientHandshake()
	if err != nil {
		fmt.Println(err)
		panic("failed handshake")
	}

	s.WriteMessage([]byte("ping")) 

}
