package entity

type Column struct {
	ID          string `json:"id"`
	BoardID     string `json:"board_id"`
	Name        string `json:"name"`
	OrderNumber int    `json:"order_number"`
}

type Columns []Column

func (c Columns) GetIDs() []string {
	var ids []string
	for _, column := range c {
		ids = append(ids, column.ID)
	}
	return ids
}

type ColumnCreate struct {
	Name       string `json:"name"`
	BoardID    string `json:"board_id"`
	OderNumber int    `json:"order"`
}

type ColumnUpdate struct {
	ID          string
	Name        *string
	OrderNumber *int
}

type GetColumnsOpts struct {
	ID          string
	BoardID     string
	Name        string
	OrderNumber int
}
