package dto

type TypePOST struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
}

type TypeGetUri struct {
	TypeID string `json:"type" binding:"required"`
}

type TypePATCH struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

type TypePATCHUri struct {
	TypeID string `json:"type" binding:"required"`
}

type TypeDELETEUri struct {
	TypeID string `json:"type" binding:"required"`
}
