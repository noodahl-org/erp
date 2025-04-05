package postgres

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBClient interface {
	AutoMigrate(*interface{}) error
	Fetch(ctx context.Context, query, out interface{}, preload ...string) error
	FetchMany(ctx context.Context, query interface{}, out interface{}, preload ...string) error
	Create(ctx context.Context, query interface{}) error
	Upsert(ctx context.Context, query interface{}, omit ...string) error
	Delete(ctx context.Context, query interface{}) error
}

type dbClient struct {
	pgdb *gorm.DB
}

func NewPGDB(opts ...func(*dbClient)) DBClient {
	db := &dbClient{}
	for _, opt := range opts {
		opt(db)
	}
	return db
}

func WithGormDB(pgdb *gorm.DB) func(*dbClient) {
	return func(c *dbClient) {
		c.pgdb = pgdb
	}
}

func (db *dbClient) AutoMigrate(model *interface{}) error {
	return db.pgdb.AutoMigrate(model)
}

func (db *dbClient) FetchMany(ctx context.Context, query interface{}, out interface{}, preload ...string) error {
	tx := db.pgdb.WithContext(ctx).
		Where(query)
	for _, pre := range preload {
		tx.Preload(pre)
	}
	tx.Find(out)
	return tx.Error
}

func (db *dbClient) Fetch(ctx context.Context, query, out interface{}, preload ...string) error {
	tx := db.pgdb.WithContext(ctx).
		Limit(1).Where(query)
	for _, pre := range preload {
		tx.Preload(pre)
	}
	tx.Find(out)
	return tx.Error
}

func (db *dbClient) Create(ctx context.Context, query interface{}) error {
	tx := db.pgdb.WithContext(ctx)
	tx.Create(query)
	return tx.Error
}

func (db *dbClient) Upsert(ctx context.Context, query interface{}, omit ...string) error {
	tx := db.pgdb.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		UpdateAll: true,
	})
	for _, o := range omit {
		tx.Omit(o)
	}
	tx.Save(query)
	return tx.Error
}

func (db *dbClient) Delete(ctx context.Context, query interface{}) error {
	tx := db.pgdb.Delete(query)
	return tx.Error
}
