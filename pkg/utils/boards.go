package utils

import "github.com/sater-151/todo-list/internal/entity"

func BoardsExcept(boards []entity.Board, boardID string) entity.Boards {
	var result entity.Boards
	for _, board := range boards {
		if board.ID != boardID {
			result = append(result, board)
		}
	}
	return result
}
