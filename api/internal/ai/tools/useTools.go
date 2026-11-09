package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ArthurMilligan/charlie/api/internal/app/domain/interfaces"
)

func UseTools(ctx context.Context, n interfaces.NotesService, toolName string, args string, messageID uint32) (string, error) {
	switch toolName {
	case "add_note":
		var msg struct {
			Content string `json:"content"`
		}

		if err := json.Unmarshal([]byte(args), &msg); err != nil {
			fmt.Println("Ошибка разбора JSON:", err)
			return UnknownToolError, err
		}

		_, err := n.AddNote(ctx, msg.Content, messageID)
		if err != nil {
			return UnknownToolError, err
		}
		return AddNoteSuccess, nil
	case "get_notes":
		notes, err := n.GetNotes(ctx)
		if err != nil {
			return UnknownToolError, err
		}

		notesJSON, err := json.Marshal(notes)
		if err != nil {
			return UnknownToolError, err
		}
		return string(notesJSON), nil
	default:
		fmt.Println("Неизвестная команда")
		return UnknownToolError, nil
	}
}
