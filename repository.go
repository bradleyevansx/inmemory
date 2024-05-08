package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type IRepository[T IEntity] interface {
	Get()([]T, error)
	GetById(id string)(*T, error)
	Create(e *T)(*T, error)
	Delete(id string)(error)
	Update(e *T)(*T, error)
}

type BaseRepository[T IEntity] struct {
	db        *sql.DB
	tableName string
	rowMapper func(scan func(dest ...any) error) (*T, error)
}

func NewRepository[T IEntity](db *sql.DB, tableName string, rowMapper  func(scan func(dest ...any) error) (*T, error)) *BaseRepository[T] {
	return &BaseRepository[T]{db: db, tableName: tableName, rowMapper: rowMapper}
}

func (r *BaseRepository[T]) Get() ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	fmt.Println(query)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var results []T

	for rows.Next() {
		res, err := r.rowMapper(rows.Scan)
		if err != nil{
			return nil, fmt.Errorf("error parsing row: %w", err)
		}

		results = append(results, *res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	fmt.Printf("all %s retrieved: %v", r.tableName, results)
	return results, nil
}

func (r *BaseRepository[T]) GetById(id string) (*T, error){
	query := fmt.Sprintf("SELECT * FROM %s as e WHERE e.id = '%s'", r.tableName, id)
	fmt.Println(query)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rows.Next()
	entity, err := r.rowMapper(rows.Scan)
	if err != nil{
		return nil, fmt.Errorf("error parsing row: %w", err)	
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	fmt.Printf("%s retrieved by id: %v", r.tableName, entity)
	return entity, nil
}

func (r *BaseRepository[T]) Create(e *T) (*T, error) {
    destructuredEntity, err := destructureEntity(e)
	if err != nil {
		return nil, fmt.Errorf("error destructuring entity: %v", err)
	}
    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *;", r.tableName, strings.Join(destructuredEntity.fieldNames, ", "), strings.Join(destructuredEntity.fieldValues, ", "))
    fmt.Println(query)
    
    row := r.db.QueryRow(query)
	res, err := r.rowMapper(row.Scan)
    if err != nil {
        return nil,fmt.Errorf("error iterating over row: %w", err)
    }
    return res, nil
}

func (r *BaseRepository[T]) Delete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s as e WHERE e.id = '%s';", r.tableName, id)
	fmt.Println(query)
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("Delete: error executing delete query: %w", err)
	}
	return nil
}

func (r *BaseRepository[T]) Update(e *T) (*T, error) {
    destructuredEntity, err := destructureEntity(e)
	if err != nil {
		return nil, fmt.Errorf("error destructuring entity: %v", err)
	}
	var setters []string
	for i := 0; i < len(destructuredEntity.fieldNames); i++ {
		if destructuredEntity.fieldNames[i] == "Entity" {
			continue
		};
		setters = append(setters, fmt.Sprintf("%s = %s", destructuredEntity.fieldNames[i], destructuredEntity.fieldValues[i]))
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %s RETURNING *;", r.tableName, strings.Join(setters, ", "), destructuredEntity.id)
	fmt.Println(query)
	row := r.db.QueryRow(query)
	res, err := r.rowMapper(row.Scan)
    if err != nil {
        return nil,fmt.Errorf("error iterating over row: %w", err)
    }
    return res, nil
}

type DestructuredEntity struct {
	id string
	fieldNames []string
	fieldValues []string
}

func destructureEntity[T IEntity](e *T)(*DestructuredEntity, error){
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

