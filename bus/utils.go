package bus

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bradleyevansx/inmemory/stor"
)

type DestructuredEntity struct {
	id string
	fieldNames []string
	fieldValues []string
}

func destructureEntity[T stor.IEntity](e *T)(*DestructuredEntity, error){
	value := reflect.ValueOf(e)
    if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
        return nil, fmt.Errorf("destructureEntity: expects a pointer to a struct")
    }
	var fieldNames []string
    var fieldValues []string
	id := ""
    for i := 0; i < value.Elem().NumField(); i++ {
        fieldValue := value.Elem().Field(i)
        if !fieldValue.IsValid() {
            continue
        }

        field := value.Elem().Type().Field(i)
        if field.PkgPath != "" {
            continue
        }

        zero := reflect.Zero(field.Type)
        if reflect.DeepEqual(fieldValue.Interface(), zero.Interface()) {
            continue
        }

        jsonTagName := field.Tag.Get("json")
			
		if field.Name == "Entity" {
			id = fmt.Sprintf("%v", fieldValue.Interface())
		}
        if jsonTagName != "" {
            fieldNames = append(fieldNames, jsonTagName)
        } else {
            fieldNames = append(fieldNames, field.Name)
        }
        fieldValues = append(fieldValues, fmt.Sprintf("'%v'", fieldValue.Interface()))
    }
	return &DestructuredEntity{
		id: fmt.Sprintf("'%s'", strings.Trim(id, "{}")),
		fieldNames: fieldNames,
		fieldValues: fieldValues,
	}, nil
}