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
	"errors"
	"fmt"
)

type Message struct {
	AlarmChannel string
	Extra        string
}

const (
	PartLength = iota
	PartChannelLength
	PartChannel
	PartExtraLength
	PartExtra
)

func BytesToMessage(bytes []byte) (m Message, err error) {
	state := 0
	messageLength := int(0)
	channelLength := int(0)
	channel := ""
	extraLength := int(0)
	extra := ""
	for i := 0; i < len(bytes); i++ {
		switch state {
		case PartLength:
			messageLength = int(bytes[i])
			if messageLength > 255 {
				err = errors.New(fmt.Sprintf("Message exceeds 255 bytes, ignoring."))
			} else if len(bytes) < messageLength {
				err = errors.New(fmt.Sprintf("Transferred bytes is less than specified message length, ignoring."))
			}
			state++
		case PartChannelLength:
			channelLength = int(bytes[i])
			if channelLength > 127 {
				err = errors.New(fmt.Sprintf("Channel name exceeds 127 bytes, ignoring."))
			} else if (len(bytes) - i) < channelLength {
				err = errors.New(fmt.Sprintf("Remaining bytes are unsufficient for the specified channel length."))
			}
			state++
		case PartChannel:
			channel = string(bytes[i : i+channelLength])
			i += channelLength
			state++
		case PartExtraLength:
			extraLength = int(bytes[i])
			if extraLength == 0 {
				break
			} else if (len(bytes) - i) < extraLength {
				err = errors.New(fmt.Sprintf("Remaining bytes are unsufficient for the specified extra data length."))
			}
			state++
		case PartExtra:
			extra = string(bytes[i : i+extraLength])
			i += extraLength
			break
		default:
			err = errors.New(fmt.Sprintf("Unknown message part %d!", bytes[i]))
		}
		if err != nil {
			return
		}
	}
	m.AlarmChannel = channel
	m.Extra = extra
	return
}

func MessageToBytes(m Message) (bytes []byte, err error) {
	messageLength := len(m.AlarmChannel) + len(m.Extra) + 3
	if messageLength > 255 {
		err = errors.New(fmt.Sprintf("Message size exceeds maximum value! %d/%d", messageLength, 255))
		return
	}
	// This is probably a terrible way to do this.
	bytes = append(
		[]byte{uint8(messageLength), uint8(len(m.AlarmChannel))},
		[]byte(m.AlarmChannel)...,
	)
	bytes = append(
		bytes,
		[]byte{uint8(len(m.Extra))}...,
	)
	bytes = append(
		bytes,
		[]byte(m.Extra)...,
	)
	return
}
