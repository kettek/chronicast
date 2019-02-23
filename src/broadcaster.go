/*
Copyright 2019 Ketchetwahmeegwun T. Southall

This file is part of chronicast.

chronicast is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

chronicast is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with chronicast.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"net"
)

type Broadcaster struct {
	Connection *net.UDPConn
	Addr       *net.UDPAddr
}

func NewBroadcaster(address string) (b Broadcaster, err error) {
	b = Broadcaster{}
	b.Addr, err = net.ResolveUDPAddr("udp", address)
	return
}

func (b *Broadcaster) Open() (err error) {
	b.Connection, err = net.DialUDP("udp", nil, b.Addr)
	return
}

func (b *Broadcaster) Send(m Message) (err error) {
	bytes, err := MessageToBytes(m)
	if err != nil {
		return
	}
	_, err = b.Connection.Write(bytes)
	return
}
