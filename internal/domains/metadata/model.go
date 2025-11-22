package metadata

import (
	"agentic/commerce/internal/core"
	"database/sql/driver"
	"fmt"

	"github.com/go-json-experiment/json/v1"
)

type MetaDataModel struct {
	core.BaseModel
	UUid     *string `gorm:"Column:uuid"`
	UserId   *int64  `gorm:"Column:user_id"`
	Metadata JSONB   `gorm:"Column:metadata;type:jsonb"`
}

type JSONB map[string]interface{}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = JSONB{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid JSONB value: %T", value)
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}
