package objects

import (
	"time"
)

//Define enums for status, rating
type status uint

const (
	CHECKEDIN status = iota
	CHECKEDOUT
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
	ID string `gorm:"primary_key" json:"id,omitempty"`

	// General details
	Title       string    `json:"title,omitempty"`
	Author      string    `json:"author,omitempty"`
	PublishDate time.Time `json:"publishdate,omitempty"`
	Status      status    `json:"status,omitempty"`
	Rating      rating    `json:"rating,omitempty"`
	// rating (1-3)
	// status (checkedin, checkout)

	// Change status
	//Status BookStatus `json:"status,omitempty"`

	// Meta information
	CreatedOn time.Time `json:"created_on,omitempty"`
	UpdatedOn time.Time `json:"updated_on,omitempty"`
}
