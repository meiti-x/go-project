package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-json-experiment/json/v1"
	"github.com/labstack/echo/v4"
)

// JsonV2 implements echo.JSONSerializer using encoding/json/v2
type JsonV2 struct{}

func (j *JsonV2) Serialize(c echo.Context, i interface{}, indent string) error {
	//return jsonv2.Marshal(c.Response(), )
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

func (j *JsonV2) Deserialize(c echo.Context, i interface{}) error {
	//return jsonv2.Unmarshal(in, out, opts...)
	err := json.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}
