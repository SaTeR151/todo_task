package dto

type TaskPOST struct {
	Label       string `json:"label" binding:"required"`
	TypeID      string `json:"type_id" binding:"required"`
	ColumnID    string `json:"column_id"`
	Description string `json:"description"`
}

type TaskGETUri struct {
	TaskID string `json:"task" binding:"required"`
}

type TaskPATCH struct {
	Label       *string `json:"label"`
	TypeID      *string `json:"type_id"`
	Description *string `json:"description"`
}

type TaskPATCHUri struct {
	TaskID string `json:"task" binding:"required"`
}

type TaskDELETEUri struct {
	TaskID string `json:"task" binding:"required"`
}

type TaskMOVEUri struct {
	TaskID string `json:"task" binding:"required"`
}

type TaskMOVE struct {
	ColumnDestenation string `json:"column_id" binding:"required"`
}
