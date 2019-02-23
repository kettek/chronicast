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
  "log"
  "fmt"
  "strings"
)

type Chronicast struct {
  config Config
  commander Commander
  listener Listener
}

func (c *Chronicast) Go() {
  var err error
  if err = c.config.Load(); err != nil {
    log.Print(err)
    return
  }

  // Create CLI
  if c.commander, err = NewCommander(); err != nil {
    log.Print(err)
    return
  }
  c.commander.AddCommand("ring", "Rings a given alarm at a specific time or offset.\n\tring \"My Alarm\" +5H30M", c.Ring)
  c.commander.AddCommand("alarm", "Provides interfaces to an alarm.\n\talarm \"alarm title\" \"mplayer ~/alarm.ogg\"", c.Alarm)

  go c.commander.Listen()

  // Create UDP Listener
  if c.listener, err = NewListener(c.config.Address); err != nil {
    log.Print(err)
    return
  }
  go c.listener.Listen()
  log.Printf("Now listening on %s.\n", c.config.Address)

  for {
    quit := false
    select {
    case in := <-c.listener.In:
      log.Printf("Recvd: %s\n", in)
    case cin := <-c.commander.In:
      if cin == "quit" {
        quit = true
      }
    }
    if quit {
      break
    }
  }
}

func (c *Chronicast) Alarm(args []string) {
  if len(args) <= 2 {
    c.commander.ShowHelp([]string{"help", "alarm"})
  } else {
    fmt.Printf("Issuing alarm %s with \"%s\"\n", args[1], strings.Join(args[2:], " "))
  }
}

func (c *Chronicast) Ring(args []string) {
  if len(args) == 2 {
    fmt.Printf("Ringing alarm %s immediately\n", args[1])
  } else if len(args) == 3 {
    fmt.Printf("Ringing alarm %s in %s\n", args[1], args[2])
  } else {
    c.commander.ShowHelp([]string{"help", "ring"})
  }
}
