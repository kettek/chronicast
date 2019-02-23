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
  "log"
)

type Listener struct {
  In  chan string
  Addr *net.UDPAddr
  Connection *net.UDPConn
}

func NewListener(address string) (listener Listener, err error) {
  l := Listener{
    In: make(chan string),
  }
  if l.Addr, err = net.ResolveUDPAddr("udp", address); err != nil {
    return
  }

  return l, err
}

func (l *Listener) Listen() (err error) {
  if l.Connection, err = net.ListenMulticastUDP("udp", nil, l.Addr); err != nil {
    return
  }

  l.Connection.SetReadBuffer(1024)

  return
}

func (l *Listener) Loop() (err error) {
  for {
    buffer := make([]byte, 1024)
    bytes, src, err := l.Connection.ReadFromUDP(buffer)
    if err != nil {
      log.Fatal(err)
      break
    }
    log.Printf("%s %d %s", src, bytes, buffer)
  }
  return
}
