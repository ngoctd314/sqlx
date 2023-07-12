# sqlx 

A self taught reflect and sql driver

Reference: https://github.com/jmoiron/sqlx

## Hanlde Types

sqlx is intended to have the same feel as database/sql. There are 4 main handle types:

- sqlx.DB - analagous to sql.DB, a representation of a database
- sqlx.Tx - analagous to sql.Tx, a representation of a transaction
- sqlx.Stmt - analagous to sql.Stmt, a representation of a prepared statement
- sqlx.NamedStmt - a representation of a prepared statement with support for named parameters

## Connecting to your database

A DB instanct is not a connection, but an abstraction representing a Database. This is why creating a DB does not return an error and will not panic. It maintains a connection pool internally, and will attempt to connect when a connection is first needed.

## bindvars

A common misconception with bindvars is that they are used for interpolation. They are only for parameterization, and are not allowed to change the structure of an SQL statement.

```sql
-- doesn't work
db.Query("SELECT * FROM ?", "mytable")
-- doesn't work
db.Query("SELECT ?, ? FROM people", "name", "location")
```

You should treat the Rows like a database cursor rather than a materialized list of results. If you do not iterate over a whole rows result, be sure to call rows.Close() to return the connection back to the pool.

The connection used by the Query remains active until either all rows are exhausted by the iteration via Next, or rows.Close() is called, at which point it is released.

The sqlx extension Queryx behaves exactly as Query does, but returns an sqlx.Rows which has extended scanning behaviors:

```go
type Place struct {
    Country string
    City sql.NullString
    TelephoneCode int `db:"telcode"`
}

rows, err := db.Queryx("SELECT * FROM place")
for rows.Next() {
    var p Place
    err = rows.StructScan(&p)
}
```