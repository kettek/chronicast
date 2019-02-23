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
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Address string `json:"Address"`
}

func (c *Config) Load() (err error) {
	// First ensure defaults.
	c.Address = "239.0.0.0:19919"
	// Now load.
	file, err := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	bytes, _ := ioutil.ReadAll(file)

	json.Unmarshal([]byte(bytes), &c)
	// TODO: Some form of config sanity checking.
	return
}
