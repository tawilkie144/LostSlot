package Data

import (
	"LostSlot/src/Data/postgres"
	"LostSlot/src/Entities"
	"context"
	"errors"
)

// DataStore
/*
DataStore is a thin adapter to allow multiple databases to be used seamlessly
*/
type DataStore interface {
	NewConnection(host string, port string, database string, user string, password string) error
	//Query
	/*
		Query the database with context context, and sql string sql and returns the resulting rows.
		args should be a map of the parameters in the form of stringParamName->anyParamValue
		Supports parameters as @paramName and $1...
		Important: Do not mix parameterization types. Use only $1, $2... OR @param1 @param2... etc
	*/
	Query(context context.Context, sql string, args map[string]any) (any, error)
	//Exec
	/*
		Executes a NonQuery of the database with context context, and sql string sql and returns the number of affected rows.
		args should be a map of the parameters in the form of stringParamName->anyParamValue
		Supports parameters as @paramName and $1...
		Important: Do not mix parameterization types. Use only $1, $2... OR @param1 @param2... etc
	*/
	Exec(context context.Context, sql string, args map[string]any) (any, error)
	Reset()
	Close()
}

func NewDataStore(source string, config Entities.Config) (DataStore, error) {
	switch source {
	case "postgres":
		postgresStore, err := postgres.NewPostgresDataStore(config)
		return &postgresStore, err
	}
	return nil, errors.New("invalid source")
}

// QueryRows
/*
QueryRows is a thin adapter to allow multiple databases to be used seamlessly
*/
type QueryRows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}
