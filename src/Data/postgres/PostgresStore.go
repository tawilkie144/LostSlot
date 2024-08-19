package postgres

import (
	"LostSlot/src/Entities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
	"sync"
	"time"
)

func NewPostgresDataStore(config Entities.Config) (PostgresStore, error) {
	return PostgresStore{}, nil
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

type PostgresRow struct {
	row *pgx.Rows
}

var pgOnce sync.Once

func (pgs PostgresStore) NewConnection(host string, port string, database string, user string, password string) error {
	var anError error
	var dbPool *pgxpool.Pool
	pgOnce.Do(func() {
		dbConfigPtr, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, database))
		dbConfigPtr.MaxConns = 10
		dbConfigPtr.MinConns = 0
		dbConfigPtr.MaxConnLifetime = time.Hour
		dbConfigPtr.MaxConnIdleTime = time.Minute * 10
		dbConfigPtr.HealthCheckPeriod = time.Minute
		if err != nil {
			anError = fmt.Errorf("error generating Postgres config: %s", err)
			return
		}
		dbPool, err = pgxpool.NewWithConfig(context.Background(), dbConfigPtr) //.New(context.Background(), connectionString)
		if err != nil {
			anError = fmt.Errorf("error connecting to Postgres database: %s", err)
			return
		}
	})
	if anError != nil {
		return anError
	}
	pgs.pool = dbPool
	return nil
}

func (conn *PostgresStore) Query(context context.Context, sql string, args map[string]any) (any, error) {
	if conn.pool == nil {
		return nil, fmt.Errorf("error: invalid pool")
	}
	var cleanedArgs []any
	if argCount := strings.Count(sql, "@"); argCount > 0 {
		sql, cleanedArgs = cleanNamedArgs(sql, args, argCount)
	} else {
		cleanedArgs = make([]any, 0, len(args))
		for k, v := range args {
			kInt, err := strconv.ParseInt(k, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("error: When using numbered parameters, use $ and an integer. recieved: %s", k)
			}
			if kInt < 1 || kInt > int64(len(args)) {
				return nil, fmt.Errorf("error: numbered parameters must be in order from 1 to number of parameters. recieved: %s. Valid range: [1, %d]", k, len(args))
			}
			cleanedArgs = append(cleanedArgs, v)
		}
	}

	rows, err := conn.pool.Query(context, sql, cleanedArgs...)
	if err != nil {
		return nil, err
	}

	rRow := &PostgresRow{&rows}
	return rRow, nil
}

func (conn *PostgresStore) Exec(context context.Context, sql string, args map[string]any) (any, error) {
	if conn.pool == nil {

		return 0, fmt.Errorf("error: invalid pool")
	}

	cmdTag, err := conn.pool.Exec(context, sql, args)
	if err != nil {
		return 0, err
	}

	return cmdTag.RowsAffected(), nil
}

func (conn *PostgresStore) Reset() {
	conn.pool.Reset()
}

func (conn *PostgresStore) Close() {
	conn.pool.Close()
}

func (row *PostgresRow) Close() {
	(*row.row).Close()
}

func (row *PostgresRow) Err() error {
	return (*row.row).Err()
}

func (row *PostgresRow) Next() bool {
	return (*row.row).Next()
}

func (row *PostgresRow) Scan(dest ...any) error {
	return (*row.row).Scan(dest...)
}

func cleanNamedArgs(sql string, namedArgsMap map[string]any, argCount int) (string, []any) {
	stringsList := strings.Split(strings.Map(func(r rune) rune {
		switch {
		case r == '(' || r == ')' || r == ';' || r == ',':
			return ' '
		}
		return r
	}, sql), "@")[1:]
	rArray := make([]any, 0, argCount)
	cleanedSql := strings.Repeat(sql, 1)
	for i, words := range stringsList {
		arg := strings.Fields(words)[0]
		cleanedSql = strings.Replace(cleanedSql, "@"+arg, "$"+strconv.Itoa(i+1), 1)
		rArray[i] = namedArgsMap[arg]
	}

	return cleanedSql, rArray
}
