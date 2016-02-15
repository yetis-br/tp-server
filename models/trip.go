package models

import (
	"encoding/json"
	"log"
	"time"

	"github.com/yetis-br/tp-server/util"
)

//Trip defines a resource structure
type Trip struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Title       string    `json:"title" gorethink:"title"`
	Public      bool      `json:"private" gorethink:"private"`
	ManagerID   string    `json:"managerId" gorethink:"managerId"`
	CreatedDate time.Time `json:"createdDate" gorethink:"createdDate"`
	UpdatedDate time.Time `json:"updatedDate" gorethink:"updatedDate"`
}

//Initialize load Trip fields with default value
func (t Trip) Initialize() {
	t.CreatedDate = time.Now()
	t.UpdatedDate = time.Now()
	t.Public = false
}

//Validate check if all required fields have been filled
func (t *Trip) Validate() bool {
	if t.Title != "" {
		return true
	}
	return false
}

//ToJSON returns a JSON string with the Trip
func (t *Trip) ToJSON() string {
	tripJSON, err := json.Marshal(t)
	util.FailOnError(err, "Failed to convert Trip Struct to JSON")
	return string(tripJSON)
}

//LoadJSON load a Trip with a JSON string
func (t *Trip) LoadJSON(tripJSON string) {
	err := json.Unmarshal([]byte(tripJSON), &t)
	util.FailOnError(err, "Failed to load Trip from JSON")
}

//Log crates a new log with the trip fields
func (t *Trip) Log() {
	log.Printf("[Trip]Log: %v", t)
}
