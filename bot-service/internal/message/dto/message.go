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
	Command Command `json:"command"`
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
	UserID  string `json:"from_user_id"`
	RoomID  string `json:"room_id"`
	Message string `json:"message"`
}
