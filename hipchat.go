package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"github.com/andybons/hipchat"
)

type messageColor string

// These colors model the colors mentioned here: https://www.hipchat.com/docs/api/method/rooms/message
const (
	ColorYellow messageColor = "yellow"
	ColorRed    messageColor = "red"
	ColorGreen  messageColor = "green"
	ColorPurple messageColor = "purple"
	ColorGray   messageColor = "gray"
	ColorRandom messageColor = "random"
)

func sendMessageToRoom(token, roomID, from, message string, color messageColor) error {
	c := hipchat.Client{AuthToken: token}

	req := hipchat.MessageRequest{
		RoomId:        roomID,
		From:          from,
		Message:       message,
		Color:         string(color),
		MessageFormat: hipchat.FormatText,
		Notify:        false,
	}

	return c.PostMessage(req)
}
