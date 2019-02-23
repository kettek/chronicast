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
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type Commander struct {
	In           chan string
	Scanner      *bufio.Scanner
	Commands     map[string]func([]string)
	CmdRegexp    *regexp.Regexp
	Help         map[string]string
	ExtendedHelp map[string]map[string]string
}

func NewCommander() (c Commander, err error) {
	c = Commander{
		In:           make(chan string),
		Commands:     make(map[string]func([]string)),
		Help:         make(map[string]string),
		ExtendedHelp: make(map[string]map[string]string),
		Scanner:      bufio.NewScanner(os.Stdin),
	}

	if c.CmdRegexp, err = regexp.Compile("(\"[^\"]+\"|[^\\s\"]+)"); err != nil {
		return
	}

	c.AddCommand("help", "Shows help for a given topic.", c.ShowHelp)
	c.AddCommand("quit", "Quits the program.", nil)
	return
}

func (c *Commander) ShowHelp(args []string) {
	if len(args) == 1 {
		for k, _ := range c.Commands {
			fmt.Printf("%s - %s\n", k, c.Help[k])
		}
	} else {
		if len(args) == 2 {
			k := args[1]
			fmt.Printf("%s - %s\n", k, c.Help[k])
			if _, ok := c.ExtendedHelp[k]; ok {
				for k2, _ := range c.ExtendedHelp[k] {
					fmt.Printf("%s - %s\n", k2, c.ExtendedHelp[k][k2])
				}
			}
		} else {
			k := args[1]
			if _, ok := c.ExtendedHelp[k]; ok {
				for k2, _ := range c.ExtendedHelp[k] {
					fmt.Printf("%s - %s\n", k2, c.ExtendedHelp[k][k2])
				}
			}
		}
	}
}

func (c *Commander) AddCommand(k string, h string, f func([]string)) {
	c.Commands[k] = f
	c.Help[k] = h
}
func (c *Commander) AddExtendedHelp(k string, k2 string, h string) {
	if _, ok := c.ExtendedHelp[k]; !ok {
		c.ExtendedHelp[k] = make(map[string]string)
	}
	c.ExtendedHelp[k][k2] = h
}

func (c *Commander) Listen() {
	for c.Scanner.Scan() {
		args := c.CmdRegexp.FindAllString(c.Scanner.Text(), -1)
		if len(args) == 0 {
			c.Prompt()
		} else if _, ok := c.Commands[args[0]]; ok {
			if args[0] == "quit" {
				c.In <- "quit"
			} else {
				c.Commands[args[0]](args)
				c.Prompt()
			}
		} else {
			fmt.Printf("No such command \"%s\"\n", args[0])
			c.Prompt()
		}
	}
}

func (c *Commander) Prompt() {
	fmt.Print("> ")
}
