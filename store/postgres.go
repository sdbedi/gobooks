package store

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redeam/gobooks/errors"
	"github.com/redeam/gobooks/objects"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type pg struct {
	db *gorm.DB
}

// NewPostgresBookStore returns a postgres implementation of Book store
func NewPostgresBookStore(conn string) IBookStore {
	// create database connection
	db, err := gorm.Open(postgres.Open(conn),
		&gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "", log.LstdFlags),
				logger.Config{
					LogLevel: logger.Info,
					Colorful: true,
				},
			),
		},
	)
	if err != nil {
		panic("Enable to connect to database: " + err.Error())
	}
	if err := db.AutoMigrate(&objects.Book{}); err != nil {
		panic("Enable to migrate database: " + err.Error())
	}
	// return store implementation
	return &pg{db: db}
}

func (p *pg) Get(ctx context.Context, in *objects.GetRequest) (*objects.Book, error) {
	bk := &objects.Book{}
	// take book where id == uid from database
	err := p.db.WithContext(ctx).Take(bk, "id = ?", in.ID).Error
	if err == gorm.ErrRecordNotFound {
		// not found
		return nil, errors.ErrBookNotFound
	}
	return bk, err
}

func (p *pg) List(ctx context.Context, in *objects.ListRequest) ([]*objects.Book, error) {
	if in.Limit == 0 || in.Limit > objects.MaxListLimit {
		in.Limit = objects.MaxListLimit
	}
	query := p.db.WithContext(ctx).Limit(in.Limit)
	if in.Title != "" {
		query = query.Where("title ilike ?", "%"+in.Title+"%")
	}
	list := make([]*objects.Book, 0, in.Limit)
	fmt.Println(list)
	err := query.Order("id").Find(&list).Error
	return list, err
}

func (p *pg) Create(ctx context.Context, in *objects.CreateRequest) error {
	if in.Book == nil {
		return errors.ErrObjectIsRequired
	}
	in.Book.ID = GenerateUniqueID()

	in.Book.CreatedOn = p.db.NowFunc()
	return p.db.WithContext(ctx).
		Create(in.Book).
		Error
}

func (p *pg) UpdateDetails(ctx context.Context, in *objects.UpdateDetailsRequest) error {
	bk := &objects.Book{
		ID:          in.ID,
		Title:       in.Title,
		Author:      in.Author,
		PublishDate: in.PublishDate,
		Publisher:   in.Publisher,
		Status:      in.Status,
		Rating:      in.Rating,
		UpdatedOn:   p.db.NowFunc(),
	}
	return p.db.WithContext(ctx).Model(bk).
		Select("title", "author", "publishdate", "status", "rating", "updated_on").
		Updates(bk).
		Error
}

func (p *pg) Delete(ctx context.Context, in *objects.DeleteRequest) error {
	bk := &objects.Book{ID: in.ID}
	return p.db.WithContext(ctx).Model(bk).
		Delete(bk).
		Error
}
