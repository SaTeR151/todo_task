package dto

type ColumnPOST struct {
	Name        string `json:"name" binding:"required"`
	OrderNumber int    `json:"order_number" binding:"required"`
}

type ColumnGETUri struct {
	ColumnID string `json:"column" binding:"required"`
}

type ColumnDELETEUri struct {
	ColumnID string `json:"column" binding:"required"`
}

type ColumnPATCH struct {
	Name        *string `json:"name"`
	OrderNumber *int    `json:"order_number"`
}

type ColumnPATCHUri struct {
	ColumnID string `json:"column" binding:"required"`
}

type ColumnSWAP struct {
	ColumnIDA string `json:"column_a" binding:"required"`
	ColumnIDB string `json:"column_b" binding:"required"`
}
