package models

//Trip defines a resource structure
type Trip struct {
	ID        string `json:"id" gorethink:"id,omitempty"`
	Title     string `json:"title" gorethink:"title"`
	Public    bool   `json:"private" gorethink:"private"`
	ManagerID string `json:"managerId" gorethink:"managerId"`
}
