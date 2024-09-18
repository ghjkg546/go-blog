package models

type MenuOption struct {
	Value    uint         `json:"value"`
	ID       uint         `json:"id"`
	ParentId uint         `json:"parent_id"`
	Label    string       `json:"label"`
	Children []MenuOption `json:"children,omitempty"`
}
