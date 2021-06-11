package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/redeam/gobooks/errors"
	"github.com/redeam/gobooks/objects"
	"github.com/redeam/gobooks/store"
)

// IBookHandler is implement all the handlers
type IBookHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	store store.IBookStore
}

// NewBookHandler return current IBookHandler implementation
func NewBookHandler(store store.IBookStore) IBookHandler {
	return &handler{store: store}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	if id == "" && title == "" {
		WriteError(w, errors.ErrValidBookIdIsRequired)
		return
	}
	bk, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.BookResponseWrapper{Book: bk})
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	// title
	title := values.Get("title")
	// limit
	limit, err := IntFromString(w, values.Get("limit"))
	if err != nil {
		return
	}
	// list books
	list, err := h.store.List(r.Context(), &objects.ListRequest{
		Limit: limit,
		Title: title,
	})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.BookResponseWrapper{Books: list})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, errors.ErrUnprocessableEntity)
		return
	}
	bk := &objects.Book{}
	if Unmarshal(w, data, bk) != nil {
		return
	}
	//Make sure we have a title and author
	if bk.Title == "" || bk.Author == "" {
		WriteError(w, errors.ErrTitleandAuthorIsRequired)
		return
	}
	//Check the status if we have an appropriate status - set to CheckedIn if empty, return error if a non-acceptable status is submitted
	if bk.Status != "CheckedIn" && bk.Status != "CheckedOut" {
		if bk.Status == "" {
			bk.Status = "CheckedIn"
		} else {
			WriteError(w, errors.ErrStatusIsRequired)
			return
		}
	}
	//Check that rating is supplied
	if bk.Rating > 3 || bk.Rating < 1 {
		WriteError(w, errors.ErrRatingIsRequired)
		return
	}
	if err = h.store.Create(r.Context(), &objects.CreateRequest{Book: bk}); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.BookResponseWrapper{Book: bk})
}

func (h *handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, errors.ErrUnprocessableEntity)
		return
	}
	req := &objects.UpdateDetailsRequest{}
	if Unmarshal(w, data, req) != nil {
		return
	}
	//Check if ID is supplied
	if req.ID == "" {
		WriteError(w, errors.ErrValidBookIdIsRequired)
		return
	}
	//Check the status
	if req.Status != "CheckedIn" && req.Status != "CheckedOut" && len(req.Status) > 0 {
		WriteError(w, errors.ErrStatusIsRequired)
		return
	}
	//check if book exists.
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: req.ID}); err != nil {
		WriteError(w, err)
		return
	}

	//TODO: restructure this method to return the new book object so the retrieve call below is unecessary
	if err = h.store.UpdateDetails(r.Context(), req); err != nil {
		WriteError(w, err)
		return
	}

	//Retrieve the new book ()
	bk, err := h.store.Get(r.Context(), &objects.GetRequest{ID: req.ID})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.BookResponseWrapper{Book: bk})
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteError(w, errors.ErrValidBookIdIsRequired)
		return
	}

	// check if book exist
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id}); err != nil {
		WriteError(w, err)
		return
	}

	if err := h.store.Delete(r.Context(), &objects.DeleteRequest{ID: id}); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.BookResponseWrapper{})
}
