package objects

import (
	"encoding/json"
	"net/http"
)

// MaxListLimit maximum listting
const MaxListLimit = 200

// GetRequest for retrieving single Book
type GetRequest struct {
	ID string `json:"id"`
}

// ListRequest for retrieving list of Books
type ListRequest struct {
	Limit int `json:"limit"`
	// optional title matching
	Title string `json:"title"`
}

// CreateRequest for creating a new Book
type CreateRequest struct {
	Book *Book `json:"book"`
}

// UpdateDetailsRequest to update existing Book
type UpdateDetailsRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publishdate"`
	Status      status `json:"status"`
	Rating      rating `json:"rating"`
}

// DeleteRequest to delete a Book
type DeleteRequest struct {
	ID string `json:"id"`
}

// BookResponseWrapper reponse of any Book request
type BookResponseWrapper struct {
	Book  *Book   `json:"book,omitempty"`
	Books []*Book `json:"books,omitempty"`
	Code  int     `json:"-"`
}

// JSON convert BookResponseWrapper in json
func (e *BookResponseWrapper) JSON() []byte {
	if e == nil {
		return []byte("{}")
	}
	res, _ := json.Marshal(e)
	return res
}

// StatusCode return status code
func (e *BookResponseWrapper) StatusCode() int {
	if e == nil || e.Code == 0 {
		return http.StatusOK
	}
	return e.Code
}
