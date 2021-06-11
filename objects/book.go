package objects

import (
	"time"
)

//Define enums for status, rating
type status string

const (
	CheckedIn  status = "CheckedIn"
	CheckedOut status = "CheckedOut"
)

type rating uint

const (
	R1 rating = iota + 1
	R2
	R3
)

// Book object for the API
type Book struct {
	// Identifier
	//TODO: Change to pointers for better handling of nil values
	ID string `gorm:"primary_key" json:"id,omitempty"`

	// General details
	Title     string `json:"title,omitempty"`
	Author    string `json:"author,omitempty"`
	Publisher string `json:"publisher,omitempty"`
	//TODO: implement date in custom type/struct
	PublishDate string `json:"publishdate,omitempty"`
	Status      status `json:"status,omitempty"`
	Rating      rating `json:"rating,omitempty"`

	// Meta information
	CreatedOn time.Time `json:"created_on,omitempty"`
	UpdatedOn time.Time `json:"updated_on,omitempty"`
}
