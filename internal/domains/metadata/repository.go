package metadata

import (
	"context"

	"agentic/commerce/internal/core"

	"agentic/commerce/internal/infrastructure/database"
)

type IContentRepository interface {
	core.IBaseRepository[MetaDataModel]
	ListByUserID(ctx context.Context, userId int64) ([]MetaDataModel, error)
	GetByUserID(ctx context.Context, uuid string, userId int64) (*MetaDataModel, error)
}

type contentRepository struct {
	core.IBaseRepository[MetaDataModel]
	database database.GormDB
}

func NewContentRepository(database database.GormDB) IContentRepository {
	return &contentRepository{
		IBaseRepository: core.NewBaseRepository[MetaDataModel](database),
		database:        database,
	}
}

func (db *contentRepository) ListByUserID(ctx context.Context, userId int64) ([]MetaDataModel, error) {
	var result []MetaDataModel
	tx := db.database(ctx).Model(&MetaDataModel{}).
		Where("user_id = ?", userId).
		Find(&result)
	return core.ResolveDBSliceResult(result, tx)
}

func (db *contentRepository) GetByUserID(ctx context.Context, uuid string, userId int64) (*MetaDataModel, error) {
	var result *MetaDataModel
	tx := db.database(ctx).Model(&MetaDataModel{}).
		Where("user_id = ? and uuid = ?", userId, uuid).
		First(&result)
	return core.ResolveDBResult(result, tx)
}
