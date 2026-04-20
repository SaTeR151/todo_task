package entity

type Type struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

type Types []Type

func (t Types) GetIDs() []string {
	var ids []string
	for _, task := range t {
		ids = append(ids, task.ID)
	}
	return ids
}

type GetTypesOpts struct {
	ID     string
	UserID string
}

type TypeCreate struct {
	UserID string `db:"user_id"`
	Name   string `db:"name"`
	Color  string `db:"color"`
}

type TypeUpdate struct {
	ID    string
	Name  *string
	Color *string
}
