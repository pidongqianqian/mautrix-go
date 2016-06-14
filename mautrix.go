package mautrix

import (
	"fmt"
)

// Session .
type Session struct {
	NextBatch   string
	Rooms       map[string]RoomInfo // The rooms a user is part of
	AccessToken string
	HomeServer  string
	Timeline    chan Event
	OnJoin      chan string // When we find a new room
	stop        chan bool   // stop the service
}

// Start ..
func (session *Session) Start() {
Loop:
	for {
		select {
		case <-session.stop:
			break Loop
		default:
			err := session.Sync()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// Close closes everything down
func (session *Session) Close() {
	session.stop <- true
}

// GetURL gets the URL for the given path on this session.
func (session *Session) GetURL(path string, args ...interface{}) string {
	return fmt.Sprintf("%s/_matrix/client/r0%s", session.HomeServer, fmt.Sprintf(path, args...))
}

// Init .
func Init(homeserver string) *Session {
	session := Session{HomeServer: homeserver,
		NextBatch: "s9_13_0_1_1_1",
		Timeline:  make(chan Event, 10),
		OnJoin:    make(chan string, 10),
		Rooms:     make(map[string]RoomInfo),
		stop:      make(chan bool),
	}

	return &session
}