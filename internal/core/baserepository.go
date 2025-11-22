package core

import (
	"context"

	"agentic/commerce/internal/infrastructure/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IBaseRepository[TModel any] interface {
	Save(ctx context.Context, model *TModel) error
	Create(ctx context.Context, model *TModel) error
	Update(ctx context.Context, model *TModel) error
	Upsert(ctx context.Context, model *TModel) error
	Delete(ctx context.Context, model *TModel) error
}
type baseRepository[TModel any] struct {
	database database.GormDB
}

func NewBaseRepository[TModel any](database database.GormDB) IBaseRepository[TModel] {
	return &baseRepository[TModel]{
		database: database,
	}
}

func (db *baseRepository[TModel]) Save(ctx context.Context, model *TModel) error {
	tx := db.database(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(model)
	return ResolveError(tx)
}

func (db *baseRepository[TModel]) Create(ctx context.Context, model *TModel) error {
	tx := db.database(ctx).
		Create(model)
	return tx.Error
}

func (db *baseRepository[TModel]) Update(ctx context.Context, model *TModel) error {
	tx := db.database(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&model)
	return tx.Error
}

func (db *baseRepository[TModel]) Upsert(ctx context.Context, model *TModel) error {
	tx := db.database(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(model)
	return tx.Error
}

func (db *baseRepository[TModel]) Delete(ctx context.Context, model *TModel) error {
	tx := db.database(ctx).Delete(&model)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
