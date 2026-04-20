package dto

type BoardPOST struct {
	Name string `json:"name" binding:"required"`
}

type BoadGETUri struct {
	BoardID string `json:"board" binding:"required"`
}

type BoardPATCH struct {
	Name *string `json:"name"`
}

type BoardPATCHUri struct {
	BoardID string `json:"board" binding:"required"`
}

type BoardDELETEUri struct {
	BoardID string `json:"board" binding:"required"`
}
