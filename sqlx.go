package sqlx

import (
	"database/sql"
	"errors"
	"learn-sqlx/reflectx"
	"reflect"
	"strings"
)

// Execer is an interface used by MustExec and LoadFile
type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Queryer is an interface used by Get and Select
type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Queryx(query string, args ...interface{}) (*Rows, error)
}

// ColScanner in an interface used by MapScan and SliceScan
type ColScanner interface {
	Columns() ([]string, error)
	Scan(dest ...interface{}) error
	Err() error
}

// DB is a wrapper around *sql.DB which keeps track of the driverName upon Open
// used mostly to automatically bind named queries using the right bindVars
type DB struct {
	*sql.DB
	driverName string
	unsafe     bool
	Mapper     *reflectx.Mapper
}

// NewDb returns a new sqlx DB wrapper for a pre-existing *sql.DB.
// The driverName of the original database is required fro named query support.
func NewDb(db *sql.DB, driverName string) *DB {
	return &DB{DB: db, driverName: driverName}
}

// DriverName returns the driverName passed to the Open function for this DB.
func (db *DB) DriverName() string {
	return db.driverName
}

// Queryx queries the database and returns an *sqlx.Rows
// Any placeholder parameters are replaced with supplied args.
func (db *DB) Queryx(query string, args ...interface{}) (*Rows, error) {
	r, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return &Rows{
		Rows:   r,
		unsafe: db.unsafe,
		Mapper: db.Mapper,
	}, nil
}

// Open is the same as sql.Open, but returns an *sqlx.DB instead
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db, driverName: driverName}, nil
}

// Connect to a database and verify with a ping
func Connect(driverName, dataSourceName string) (*DB, error) {
	db, err := Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// MustConnect connect to a database and panics on error
func MustConnect(driverName, dataSourceName string) *DB {
	db, err := Connect(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	return db
}

// MustExec execs the query using e and panics if there was an error.
// Any placeholder parameters are placed with supplied args
func MustExec(e Execer, query string, args ...any) sql.Result {
	res, err := e.Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return res
}

// MustExec (panic) runs MustExec using this database
// Any placeholder parameters are replaced with supplied args
func (db *DB) MustExec(query string, args ...any) sql.Result {
	return MustExec(db, query, args...)
}

// Rows is a wrapper around sql.Rows which cache costly reflect operations
// during a looped StructScan
type Rows struct {
	*sql.Rows
	unsafe bool
	Mapper *reflectx.Mapper
	// these fields cache memory use for a rows during iteration w/ structScan
	started bool
	fields  [][]int
	values  []interface{}
}

// StructScan is like sql.Rows.Scan, but scans a single Row into a single Struct.
// Use this is and iterate over Rows manually when the memory load of Select() might be prohibitive.
// *Rows.StructScan caches the reflect work of matching up column positions to fields to avoid that overhead
// per scan, which means it is not safe to run StructScan on the same Rows instance with different struct types.
func (r *Rows) StructScan(dest interface{}) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}

	v = v.Elem()
	if !r.started {
		columns, err := r.Columns()
		if err != nil {
			return err
		}
		// map column => value
		for _, col := range columns {
			// r.values = append(r.values, v)
			for i := 0; i < v.NumField(); i++ {
				if v.Type().Field(i).Tag.Get("db") == col {
					r.values = append(r.values, v.Field(i).Addr().Interface())
				} else if strings.ToLower(v.Type().Field(i).Name) == col {
					r.values = append(r.values, v.Field(i).Addr().Interface())
				}
			}
		}

		r.started = true
	}

	if err := r.Scan(r.values...); err != nil {
		return err
	}

	return r.Err()
}

// SliceScan a row, returing a []interface{} with values similar to MapScan.
// This function is primarily intended for use when the number of columns is not known.
// Because you can can pass an []interface{} directly to Scan
// It's recommended that you do that as it will not have to allocated new slices per row.
func SliceScan(r ColScanner) ([]interface{}, error) {
	// columns, err := r.Columns()
	// if err != nil {
	// 	return []interface{}{}, err
	// }

	// values := make([]interface{}, len(columns))
	// for i := range values {
	// 	values[i] = new(interface{})
	// }

	// err = r.Scan(values...)

	// if err != nil {
	// 	return values, err
	// }

	// for i := range columns {
	// 	values[i] = *(values[i].(*interface{}))
	// }

	// return values, r.Err()
	return nil, nil
}
