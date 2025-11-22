package core

import (
	"errors"

	"gorm.io/gorm"
)

func ResolveError(db *gorm.DB) error {
	if nil == db.Error || errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return db.Error
}

func ResolveDBResult[TModel any](model *TModel, db *gorm.DB) (*TModel, error) {
	if nil == db.Error {
		return model, nil
	}
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, db.Error
}

func ResolveDBSliceResult[TModel any](models []TModel, db *gorm.DB) ([]TModel, error) {
	if nil == db.Error {
		return models, nil
	}
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return models, nil
	}
	return models, db.Error
}
