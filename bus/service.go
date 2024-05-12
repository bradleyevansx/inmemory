package bus

import (
	"fmt"
	"strings"

	"github.com/bradleyevansx/inmemory/stor"
)

type IService[T stor.IEntity] interface {
	Get()([]T, error)
	GetById(id string)(*T, error)
	Create(e *T)(*T, error)
	Delete(id string)(error)
	Update(e *T)(*T, error)
}

type BaseService[T stor.IEntity] struct {
	repo *stor.Repository
	tableName string
	rowMapper func(scan func(dest ...any) error) (*T, error)
}

func NewBaseService[T stor.IEntity](repo *stor.Repository, tableName string, rowMapper  func(scan func(dest ...any) error) (*T, error)) *BaseService[T] {
	return &BaseService[T]{repo: repo, tableName: tableName, rowMapper: rowMapper}
}

func (r *BaseService[T]) Get() ([]T, error) {
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	rows, err := r.repo.Query(query)
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
	return results, nil
}

func (r *BaseService[T]) GetById(id string) (*T, error){
	query := fmt.Sprintf("SELECT * FROM %s as e WHERE e.id = '%s'", r.tableName, id)
	rows, err := r.repo.Query(query)
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
	return entity, nil
}

func (r *BaseService[T]) Create(e *T) (*T, error) {
    destructuredEntity, err := destructureEntity(e)
	if err != nil {
		return nil, fmt.Errorf("error destructuring entity: %v", err)
	}
    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *;", r.tableName, strings.Join(destructuredEntity.fieldNames, ", "), strings.Join(destructuredEntity.fieldValues, ", "))
    
    row := r.repo.QueryRow(query)
	res, err := r.rowMapper(row.Scan)
    if err != nil {
        return nil,fmt.Errorf("error iterating over row: %w", err)
    }
    return res, nil
}

func (r *BaseService[T]) Delete(id string) error {
	query := fmt.Sprintf("DELETE FROM %s as e WHERE e.id = '%s';", r.tableName, id)
	_, err := r.repo.Exec(query)
	if err != nil {
		return fmt.Errorf("Delete: error executing delete query: %w", err)
	}
	return nil
}

func (r *BaseService[T]) Update(e *T) (*T, error) {
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
	row := r.repo.QueryRow(query)
	res, err := r.rowMapper(row.Scan)
    if err != nil {
        return nil,fmt.Errorf("error iterating over row: %w", err)
    }
    return res, nil
}



