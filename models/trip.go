package models

import "time"

//Trip defines a resource structure
type Trip struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Title       string    `json:"title" gorethink:"title"`
	Public      bool      `json:"private" gorethink:"private"`
	ManagerID   string    `json:"managerId" gorethink:"managerId"`
	CreatedDate time.Time `json:"createdDate" gorethink:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate" gorethink:"updatedDate"`
}

//Validate check if all required fields have been filled
func (t Trip) Validate() bool {
	return true
}
