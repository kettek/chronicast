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
	"fmt"
	"github.com/jinzhu/now"
	"log"
	"time"
)

type Chronicast struct {
	config      Config
	alarms      map[string][]string
	rings       map[string][]*time.Timer
	commander   Commander
	listener    Listener
	broadcaster Broadcaster
}

func (c *Chronicast) Init() {
	c.alarms = make(map[string][]string, 0)
	c.rings = make(map[string][]*time.Timer, 0)
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
	c.commander.AddCommand("ring", "Rings a given alarm at a specific time or offset.\n\tring \"My Alarm\" +5h30m\n\tring \"My Alarm\" 3:04PM\n\tring \"My Alarm\" 15:04\n\tring \"My Alarm\" 2019-02-22T15:04", c.Ring)
	c.commander.AddCommand("alarm", "Provides interfaces to an alarm.\n\talarm \"alarm title\" \"mplayer ~/alarm.ogg\"", c.Alarm)

	go c.commander.Listen()

	// Create UDP Listener
	if c.listener, err = NewListener(c.config.Address); err != nil {
		log.Print(err)
		return
	}
	if err = c.listener.Listen(); err != nil {
		log.Print(err)
		return
	}
	go c.listener.Loop()
	log.Printf("Now listening on %s.\n", c.config.Address)

	// Create UDP Broadcaster
	if c.broadcaster, err = NewBroadcaster(c.config.Address); err != nil {
		log.Print(err)
		return
	}
	if err = c.broadcaster.Open(); err != nil {
		log.Print(err)
		return
	}
	log.Printf("Now available for broadcasting on %s.\n", c.config.Address)

	for {
		quit := false
		select {
		case message := <-c.listener.In:
			if _, ok := c.alarms[message.AlarmChannel]; ok {
				log.Printf("Should run \"%s\"\n", c.alarms[message.AlarmChannel])
			}
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
		if _, ok := c.alarms[args[1]]; !ok {
			c.alarms[args[1]] = make([]string, len(args[2:]))
		}
		c.alarms[args[1]] = args[2:]
	}
}

func (c *Chronicast) Ring(args []string) {
	var err error
	if len(args) == 2 {
		c.onRing(args[1])
	} else if len(args) == 3 {
		var t time.Time
		var duration time.Duration
		if args[2][0] == '+' { // Duration
			duration, err = time.ParseDuration(args[2][1:])
		} else { // Assume that jinzhu's now can handle it.
			t, err = now.Parse(args[2])
			duration = time.Until(t)
		}
		fmt.Printf("Ringing alarm %s in %s\n", args[1], duration.String())

		if _, ok := c.rings[args[1]]; !ok {
			c.rings[args[1]] = make([]*time.Timer, 1)
		}
		// TODO: We need to remove the Timer pointers from the associated rings slice.
		c.rings[args[1]] = append(c.rings[args[1]], time.AfterFunc(duration, func() {
			c.onRing(args[1])
		}))
	} else {
		c.commander.ShowHelp([]string{"help", "ring"})
	}
	if err != nil {
		fmt.Printf("Problem during new ring: %v\n", err)
	}
}

func (c *Chronicast) onRing(alarm string) {
	c.broadcaster.Send(Message{AlarmChannel: alarm, Extra: ""})
	fmt.Printf("Ringing \"%s\"\n", alarm)
}
