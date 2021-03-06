package query

import (
	"database/sql"
	"github.com/vdongchina/ratgo/extend/db/query/parts"
)

type BaseQuery interface {
	Clone() BaseQuery
	Reset() error
	SetDb(db *sql.DB)
	SetTx(tx *sql.Tx)
	UnsetTx()
	GetTx() *sql.Tx
	Begin()    // Begin starts a transaction.
	Commit()   // Commit commits the transaction.
	Rollback() // Rollback aborts the transaction.
	QueryRow(sql string, args ...interface{}) *sql.Row
	QueryAll(sql string, args ...interface{}) (*sql.Rows, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Table(tableName string) BaseQuery
	Field(field interface{}) BaseQuery
	GetField() *parts.Field
	Where(expr string, value interface{}, linkSymbol ...string) BaseQuery
	Order(expr string) BaseQuery
	Limit(limit ...int) BaseQuery
	Values(valueMap map[string]interface{}) BaseQuery
	Set(valueMap map[string]interface{}) BaseQuery
	SetSqlType(sType string) error
	GetSqlType() string
	SetSql() BaseQuery
	GetSql() string
	FetchRow() *sql.Row
	FetchAll() (*sql.Rows, error)
	Modify() (sql.Result, error)
	GetDuration() Runtime
}

// 获取查询构造器
func GetQueryBuilder(driverName string) BaseQuery {
	switch driverName {
	case "mysql":
		return &MysqlQuery{
			Combine: Combine{
				Runtime: Runtime{},
			},
		}
	case "sql_server":
		return &SqlServerQuery{}
	}
	return nil
}
