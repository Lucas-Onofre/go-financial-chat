package dto

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCommand     = errors.New("invalid command")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidRoomID      = errors.New("invalid room ID")
	ErrUnsupportedCommand = errors.New("unsupported command")
)

type CommandMessage struct {
	UserID  string  `json:"user_id"`
	RoomID  string  `json:"room_id"`
	Command Command `json:"content"`
}

func (c *CommandMessage) Validate() error {
	if c.RoomID == "" {
		return ErrInvalidRoomID
	}

	if c.UserID == "" {
		return ErrInvalidUserID
	}

	if c.Command == "" {
		return ErrInvalidCommand
	}

	supportedCommand := false
	for _, cmd := range AvailableCommands {
		if strings.HasPrefix(c.Command.ToString(), cmd+"=") || c.Command.ToString() == cmd {
			supportedCommand = true
			break
		}
	}

	if !supportedCommand {
		return ErrUnsupportedCommand
	}

	return nil
}

type ResponseMessage struct {
	Type      string `json:"type"`
	RoomID    string `json:"room_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type MessageType string

const (
	MessageTypeBot   MessageType = "bot"
	MessageTypeError MessageType = "error"
)

func (mt MessageType) ToString() string {
	return string(mt)
}
