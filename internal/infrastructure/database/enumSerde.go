package database

import (
	"context"
	"fmt"
	"reflect"

	"github.com/orsinium-labs/enum"
	"gorm.io/gorm/schema"
)

type orsiniumEnumSerializer struct{}

//nolint:cyclop
func (s orsiniumEnumSerializer) Scan(
	ctx context.Context,
	field *schema.Field,
	dest reflect.Value,
	dbValue interface{},
) error {
	// 1) raw string handling (same as before)
	var str *string
	switch v := dbValue.(type) {
	case string:
		str = &v
	case []byte:
		sv := string(v)
		str = &sv
	case nil:
		// SQL NULL -> keep str nil
		str = nil
	default:
		return fmt.Errorf("Scan: expected string/[]byte, got %T", dbValue)
	}

	// 2) dest must be valid
	if !dest.IsValid() {
		return fmt.Errorf("Scan: destination reflect.Value is invalid")
	}

	// Work on a copy
	rv := dest

	// unwrap pointers safely: if pointer is nil, try to allocate it (if settable)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			if !rv.CanSet() {
				return fmt.Errorf("Scan: cannot allocate nil pointer for field %q", field.Name)
			}
			// allocate a new zero value for the pointer target
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
		if !rv.IsValid() {
			return fmt.Errorf("Scan: invalid element after deref for field %q", field.Name)
		}
	}

	// Identify the field value (either dest itself or named struct field)
	var fv reflect.Value
	if rv.Kind() == reflect.Struct {
		if field == nil {
			return fmt.Errorf("Scan: field metadata is nil")
		}
		if field.Name == "" {
			return fmt.Errorf("Scan: invalid empty field name")
		}

		var f reflect.Value
		for i := 0; i < rv.NumField(); i++ {
			if rv.Type().Field(i).Name == field.Name {
				f = rv.Field(i)
				break
			}
		}
		//f := rv.FieldByName(field.Name)
		if !f.IsValid() {
			// Avoid calling rv.Type() if rv is invalid (we checked above)
			return fmt.Errorf("Scan: no field %q on %s", field.Name, rv.Type())
		}
		fv = f
	} else {
		fv = rv
	}

	if !fv.CanSet() {
		return fmt.Errorf("Scan: field %q is not settable", field.Name)
	}

	// 3) Build a new instance of the field's exact type
	typ := fv.Type()
	var newInst reflect.Value
	if typ.Kind() == reflect.Ptr {
		newInst = reflect.New(typ.Elem())
	} else {
		newInst = reflect.New(typ)
	}

	// 4) set its .Value string field
	valueField := newInst.Elem().FieldByName("Value")
	if !valueField.IsValid() || valueField.Kind() != reflect.String {
		return fmt.Errorf("Scan: %s has no string field `Value`", typ)
	}
	if str != nil {
		valueField.SetString(*str)
	}

	// 5) assign back into target field
	if typ.Kind() == reflect.Ptr {
		if !fv.CanSet() {
			return fmt.Errorf("Scan: cannot set pointer field %q", field.Name)
		}
		fv.Set(newInst)
	} else {
		if !fv.CanSet() {
			return fmt.Errorf("Scan: cannot set field %q", field.Name)
		}
		fv.Set(newInst.Elem())
	}
	return nil
}

func (s orsiniumEnumSerializer) Value(
	ctx context.Context,
	field *schema.Field,
	dest reflect.Value,
	fieldValue interface{},
) (interface{}, error) {
	// if caller passed nil interface => treat as SQL NULL
	if fieldValue == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(fieldValue)
	if !rv.IsValid() {
		return nil, nil
	}

	// unwrap pointer values, but if pointer itself is nil -> NULL
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}

	// convert alias → enum.Member[string] if needed
	targetType := reflect.TypeOf(enum.Member[string]{})
	if rv.Type() != targetType {
		if rv.Type().ConvertibleTo(targetType) {
			rv = rv.Convert(targetType)
		} else {
			return nil, fmt.Errorf("Value: cannot convert %s to %s", rv.Type(), targetType)
		}
	}

	// now safe to type assert
	member, ok := rv.Interface().(enum.Member[string])
	if !ok {
		return nil, fmt.Errorf("expected enum.Member[string], got %T", rv.Interface())
	}
	return member.Value, nil
}
