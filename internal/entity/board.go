package entity

type Board struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

type Boards []Board

type BoardCreate struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type BoardUpdate struct {
	ID   string `json:"id"`
	Name *string
}

type GetBoardsOpts struct {
	ID     string
	UserID string
}
