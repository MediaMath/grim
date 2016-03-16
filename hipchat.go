package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

func sendMessageToRoom(token, roomID, from, message string, color string) error {
	c := hipchat.Client{AuthToken: token}

	req := hipchat.MessageRequest{
		RoomId:        roomID,
		From:          from,
		Message:       message,
		Color:         color,
		MessageFormat: hipchat.FormatText,
		Notify:        false,
	}

	return c.PostMessage(req)
}

func sendMessageToRoom2(token, roomID, from, message, color string) (err error) {
	var (
		sane    = sanitizeHipchatMessage(message)
		url     = fmt.Sprintf("https://api.hipchat.com/v2/room/%v/notification?auth_token=%v", roomID, token)
		payload = fmt.Sprintf(`{"message_format": "text", "from": "%v", "message": "%v", "color": "%v"}`, from, sane, color)
		resp    *http.Response
	)

	resp, err = http.Post(url, "application/json", bytes.NewBuffer([]byte(payload)))
	if err == nil {
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
	}

	if resp.StatusCode != http.StatusNoContent {
		err = fmt.Errorf("failed to send message with response code %v", resp.StatusCode)
	}

	return
}

func sanitizeHipchatMessage(message string) string {
	r := strings.NewReplacer(
		"\b", `\b`,
		"\f", `\f`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
		"\"", `\"`,
		"\\", `\\`,
	)

	sane := r.Replace(message)

	if len(sane) > 10000 {
		sane = sane[:10000]
	}

	return sane
}
