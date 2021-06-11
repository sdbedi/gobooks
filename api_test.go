package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/redeam/gobooks/errors"
	"github.com/redeam/gobooks/handlers"
	"github.com/redeam/gobooks/objects"
	"github.com/redeam/gobooks/store"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router    *mux.Router
	flushAll  func(t *testing.T)
	createOne func(t *testing.T, title string) *objects.Book
	getOne    func(t *testing.T, id string, wantErr bool) *objects.Book
)

func TestMain(t *testing.M) {
	log.Println("Registering")

	conn := "postgres://user:password@localhost:5432/db?sslmode=disable"
	if c := os.Getenv("DB_CONN"); c != "" {
		conn = c
	}

	router = mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	st := store.NewPostgresBookStore(conn)
	hnd := handlers.NewBookHandler(st)
	RegisterAllRoutes(router, hnd)

	flushAll = func(t *testing.T) {
		db, err := gorm.Open(postgres.Open(conn), nil)
		if err != nil {
			t.Fatal(err)
		}
		db.Delete(&objects.Book{}, "1=1")
	}

	createOne = func(t *testing.T, title string) *objects.Book {
		bk := &objects.Book{
			Title:       title,
			Author:      "Author of " + title,
			Publisher:   "Publisher of " + title,
			PublishDate: "Date of " + title,
			Status:      "CheckedIn",
			Rating:      1,
		}
		err := st.Create(context.TODO(), &objects.CreateRequest{Book: bk})
		if err != nil {
			t.Fatal(err)
		}
		return bk
	}
	getOne = func(t *testing.T, id string, wantErr bool) *objects.Book {
		bk, err := st.Get(context.TODO(), &objects.GetRequest{ID: id})
		if err != nil && wantErr {
			t.Fatal(err)
		}
		return bk
	}

	log.Println("Starting")
	os.Exit(t.Run())
}

func Do(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestUnknownEndpoints(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) *http.Request
	}{
		{
			name: "root",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
		},
		{
			name: "api-root",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
		},
		{
			name: "random",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/random", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			_ = assert.Equal(t, http.StatusNotFound, w.Code) &&
				assert.Equal(t, "404 page not found\n", string(w.Body.Bytes()))
		})
	}
}

func TestGetEndpoint(t *testing.T) {
	flushAll(t)
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) *http.Request
		message string
	}{
		{
			name: "OK",
			setup: func(t *testing.T) *http.Request {
				bk := createOne(t, "Ok")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/books?id="+bk.ID, nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			message: "",
		},
		{
			name: "NotFound",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1/books?id=32", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    errors.ErrBookNotFound.Code,
			message: errors.ErrBookNotFound.Message,
		},
		{
			name: "WithoutID",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1/books", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    errors.ErrValidBookIdIsRequired.Code,
			message: errors.ErrValidBookIdIsRequired.Message,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			assert.Equal(t, tt.code, w.Code)
			got := &objects.BookResponseWrapper{}
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
		})
	}
}

func TestListEndpoint(t *testing.T) {
	flushAll(t)
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) *http.Request
		listLen int
	}{
		{
			name: "Zero",
			setup: func(t *testing.T) *http.Request {
				flushAll(t)
				req, err := http.NewRequest(http.MethodGet, "/api/v1/books/list", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 0,
		},
		{
			name: "All",
			setup: func(t *testing.T) *http.Request {
				_ = createOne(t, "One")
				_ = createOne(t, "Two")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/books/list", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 2,
		},
		{
			name: "Limited",
			setup: func(t *testing.T) *http.Request {
				_ = createOne(t, "Three")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/books/list?limit=2", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 2,
		},
		{
			name: "Name",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1/books/list?title=e", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			got := &objects.BookResponseWrapper{}
			assert.Equal(t, tt.code, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
			assert.Equal(t, tt.listLen, len(got.Books))
		})
	}
}

func TestCreateEndpoint(t *testing.T) {
	flushAll(t)
	tests := []struct {
		name    string
		message string
		code    int
		bk      *objects.Book
	}{
		{
			name:    "Ok",
			message: "",
			code:    http.StatusOK,
			bk: &objects.Book{
				Title:     "Title",
				Author:    "Author",
				Publisher: "Publisher",
				Rating:    1,
			},
		},
		{
			name:    "Bad Status",
			message: errors.ErrStatusIsRequired.Message,
			code:    errors.ErrStatusIsRequired.Code,
			bk: &objects.Book{
				Title:     "Bad Status",
				Author:    "Author of Bad Status",
				Publisher: "Publisher of Bad Status",
				Rating:    1,
				Status:    "argle",
			},
		},

		{
			name:    "Missing Author",
			message: errors.ErrTitleandAuthorIsRequired.Message,
			code:    errors.ErrTitleandAuthorIsRequired.Code,
			bk: &objects.Book{
				Title:     "Missing Author",
				Author:    "",
				Publisher: "Publisher of Missing Author",
				Rating:    1,
			},
		},
		{
			name:    "Missing Title",
			message: errors.ErrTitleandAuthorIsRequired.Message,
			code:    errors.ErrTitleandAuthorIsRequired.Code,
			bk: &objects.Book{
				Title:     "",
				Author:    "...",
				Publisher: "Publisher of Missing Title",
				Rating:    1,
			},
		},
		{
			name:    "Bad Rating",
			message: errors.ErrRatingIsRequired.Message,
			code:    errors.ErrRatingIsRequired.Code,
			bk: &objects.Book{
				Title:     "Bad Rating",
				Author:    "Author of Bad Rating",
				Publisher: "Publisher of Bad Ratigin",
				Rating:    4,
			},
		},
		{
			name:    "No input",
			message: errors.ErrObjectIsRequired.Message,
			code:    errors.ErrObjectIsRequired.Code,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.bk)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewReader(b))
			if err != nil {
				t.Fatal(err)
			}
			w := Do(req)
			got, gotErr := &objects.BookResponseWrapper{}, &errors.Error{}
			assert.Equal(t, tt.code, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), gotErr))
			assert.Equal(t, tt.message, gotErr.Message)
			if tt.code == http.StatusOK {
				ok := assert.NotNil(t, got.Book) &&
					assert.NotEmpty(t, got.Book.ID) &&
					assert.NotEmpty(t, got.Book.CreatedOn) &&
					//Check that status defaults to CheckedIn
					assert.Equal(t, objects.CheckedIn, got.Book.Status)
				if ok {
					tt.bk.ID = got.Book.ID
					tt.bk.CreatedOn = got.Book.CreatedOn
					tt.bk.Status = got.Book.Status
					assert.Equal(t, tt.bk, got.Book)
				}
			}
		})
	}
}

func TestUpdateDetailsEndpoint(t *testing.T) {
	flushAll(t)
	reqFn := func(t *testing.T, bk *objects.Book) (*http.Request, *objects.Book) {
		var (
			b   []byte
			err error
		)
		if bk != nil {
			b, err = json.Marshal(&objects.UpdateDetailsRequest{
				ID:          bk.ID,
				Title:       bk.Title,
				Author:      bk.Author,
				Publisher:   bk.Publisher,
				PublishDate: bk.PublishDate,
				Status:      bk.Status,
				Rating:      bk.Rating,
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		req, err := http.NewRequest(http.MethodPut, "/api/v1/books/update", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		return req, bk
	}
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) (*http.Request, *objects.Book)
		message string
	}{
		{
			name: "OK",
			setup: func(t *testing.T) (*http.Request, *objects.Book) {
				bk := createOne(t, "Ok")
				bk.Author = "a"
				bk.Title = "b"
				return reqFn(t, bk)
			},
			code: http.StatusOK,
		},
		{
			name: "NotFound",
			setup: func(t *testing.T) (*http.Request, *objects.Book) {
				bk := createOne(t, "Ok")
				bk.ID = "1"
				return reqFn(t, bk)
			},
			message: errors.ErrBookNotFound.Message,
			code:    errors.ErrBookNotFound.Code,
		},
		{
			name: "Not ID",
			setup: func(t *testing.T) (*http.Request, *objects.Book) {
				bk := createOne(t, "Ok")
				bk.ID = ""
				return reqFn(t, bk)
			},
			message: errors.ErrValidBookIdIsRequired.Message,
			code:    errors.ErrValidBookIdIsRequired.Code,
		},
		{
			name: "No input",
			setup: func(t *testing.T) (*http.Request, *objects.Book) {
				return reqFn(t, nil)
			},
			message: errors.ErrObjectIsRequired.Message,
			code:    errors.ErrObjectIsRequired.Code,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, exp := tt.setup(t)
			w := Do(req)
			assert.Equal(t, tt.code, w.Code)
			if tt.message != "" {
				got := &errors.Error{}
				assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
				assert.Equal(t, tt.message, got.Message)
			} else if exp != nil {
				bk := getOne(t, exp.ID, true)
				assert.Equal(t, exp.Author, bk.Author)
				assert.Equal(t, exp.Title, bk.Title)
				assert.Equal(t, exp.Publisher, bk.Publisher)
				assert.Equal(t, exp.PublishDate, bk.PublishDate)
				assert.Equal(t, exp.Rating, bk.Rating)
			}
		})
	}
}

func TestDeleteEndpoint(t *testing.T) {
	flushAll(t)
	reqFn := func(t *testing.T, in *objects.DeleteRequest) (*http.Request, string) {
		id := ""
		if in != nil {
			id = in.ID
		}
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/books?id="+id, nil)
		if err != nil {
			t.Fatal(err)
		}
		return req, id
	}
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) (*http.Request, string)
		message string
	}{
		{
			name: "OK",
			setup: func(t *testing.T) (*http.Request, string) {
				bk := createOne(t, "Ok")
				return reqFn(t, &objects.DeleteRequest{ID: bk.ID})
			},
			code: http.StatusOK,
		},
		{
			name: "No input",
			setup: func(t *testing.T) (*http.Request, string) {
				return reqFn(t, nil)
			},
			message: errors.ErrValidBookIdIsRequired.Message,
			code:    errors.ErrValidBookIdIsRequired.Code,
		},
		{
			name: "NotFound",
			setup: func(t *testing.T) (*http.Request, string) {
				return reqFn(t, &objects.DeleteRequest{ID: "fake"})
			},
			message: errors.ErrBookNotFound.Message,
			code:    errors.ErrBookNotFound.Code,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, id := tt.setup(t)
			w := Do(req)
			assert.Equal(t, tt.code, w.Code)
			if tt.message != "" {
				got := &errors.Error{}
				assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
				assert.Equal(t, tt.message, got.Message)
			} else if id != "" {
				assert.Nil(t, getOne(t, id, false))
			}
		})
	}
}
