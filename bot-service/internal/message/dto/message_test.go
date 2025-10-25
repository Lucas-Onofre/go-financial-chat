package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandMessage_Validate(t *testing.T) {
	type want struct {
		err error
	}

	tests := []struct {
		name       string
		commandMsg CommandMessage
		want       want
	}{
		{
			name: "Given valid command message, When Validate is called, Then it should return no error",
			commandMsg: CommandMessage{
				UserID:  "user1",
				RoomID:  "room1",
				Command: Command("/stock=AAPL"),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Given command message with empty RoomID, When Validate is called, Then it should return ErrInvalidRoomID",
			commandMsg: CommandMessage{
				UserID:  "user1",
				RoomID:  "",
				Command: Command("/stock=AAPL"),
			},
			want: want{
				err: ErrInvalidRoomID,
			},
		},
		{
			name: "Given command message with empty UserID, When Validate is called, Then it should return ErrInvalidUserID",
			commandMsg: CommandMessage{
				UserID:  "",
				RoomID:  "room1",
				Command: Command("/stock=AAPL"),
			},
			want: want{
				err: ErrInvalidUserID,
			},
		},
		{
			name: "Given command message with empty Command, When Validate is called, Then it should return ErrInvalidCommand",
			commandMsg: CommandMessage{
				UserID:  "user1",
				RoomID:  "room1",
				Command: Command(""),
			},
			want: want{
				err: ErrInvalidCommand,
			},
		},
		{
			name: "Given command message with unsupported Command, When Validate is called, Then it should return ErrUnsupportedCommand",
			commandMsg: CommandMessage{
				UserID:  "user1",
				RoomID:  "room1",
				Command: Command("/unknown=AAPL"),
			},
			want: want{
				err: ErrUnsupportedCommand,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.commandMsg.Validate()
			assert.Equal(t, tt.want.err, err)
		})
	}
}
