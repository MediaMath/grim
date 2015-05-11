package grim

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"testing"
)

// HC_AUTH_TOKEN="" HC_ROOM_ID="" go test -v -run TestSendMessageToRoomSucceeds
func TestSendMessageToRoomSucceeds(t *testing.T) {
	token := getEnvOrSkip(t, "HC_AUTH_TOKEN")
	roomID := getEnvOrSkip(t, "HC_ROOM_ID")
	from := "Grim"
	message := "This is a test message."
	color := ColorRandom

	err := sendMessageToRoom(token, roomID, from, message, color)
	if err != nil {
		t.Fatal(err)
	}
}
