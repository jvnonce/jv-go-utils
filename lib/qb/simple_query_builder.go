package qb

import (
	"database/sql"
	"strconv"
	"strings"

	jve "github.com/jvnonce/jv-go-utils/lib/errors"
	jvm "github.com/jvnonce/jv-go-utils/lib/maps"
)

const (
	selectAction = "SELECT"
	updateAction = "UPDATE"
	insertAction = "INSERT"
	deleteAction = "DELETE"
)

type builder struct {
	db          *sql.DB
	tableName   string
	tableAlias  string
	action      string
	columns     []string
	params      []any
	where       string
	joins       []string
	orderBy     []string
	groupBy     string
	having      string
	limit       int
	offset      int
	sql         string
	isManualSQL bool
}

// Simple query builder interface for PostgreSQL
type QueryBuilder interface {
	// Prepares plain sql
	SQL(sql string, args ...any) QueryBuilder

	// Alias for table
	//
	// Ex.: qb.Select("users").Alias("u")
	Alias(alias string) QueryBuilder

	// Prepares select query
	//
	// Ex.: qb.Select("users").Alias("u")
	Select(table string) QueryBuilder

	// Prepares insert query
	//
	// Ex.: qb.Insert("users").Columns("name", "email").Parameters("jv", "jv19841202@gmail.com")
	Insert(table string) QueryBuilder

	// Prepares update query
	//
	// Ex.: qb.Update("users").Columns("name", "email").Parameters("jv", "jv19841202@gmail.com").Where("name=?", "jv")
	Update(table string) QueryBuilder

	// Prepares delete query
	//
	// Ex.: qb.Delete("users").Where("name=?", "jv")
	Delete(table string) QueryBuilder

	// Enumerates columns for select, update or insert queries
	//
	// Ex.: qb.Update("users").Columns("name", "email").Parameters("jv", "jv19841202@gmail.com").Where("name=?", "jv")
	Columns(columns ...string) QueryBuilder

	// Enumerates parameters for select, update or insert queries
	//
	// Ex.: qb.Update("users").Columns("name", "email").Parameters("jv", "jv19841202@gmail.com").Where("name=?", "jv")
	Parameters(params ...any) QueryBuilder

	// Alternative  for select, update or insert queries
	//
	// Ex.: qb.Update("users").ColsWithParams(jvm.M{"name": "jv", "email": "jv19841202@gmail.com"}.Where("name=?", "jv")
	ColsWithParams(in jvm.M) QueryBuilder

	// Where string with parameters
	//
	// Ex.: qb.Update("users").Columns("name", "email").Parameters("jv", "jv19841202@gmail.com").Where("name=?", "jv")
	Where(where string, args ...any) QueryBuilder

	// Join query to the select query
	//
	// Ex.: qb.Select("users").Alias("u").Join("INNER", "profile", "p", "u.id=p.user_id")
	Join(join string, tableName string, aliasName string, condition string) QueryBuilder

	// Inner join query table to the select query
	//
	// Ex.: qb.Select("users").Alias("u").InnerJoin("profile", "p", "u.id=p.user_id")
	InnerJoin(tableName string, aliasName string, condition string) QueryBuilder

	// Left join query table to the select query
	//
	// Ex.: qb.Select("users").Alias("u").LeftJoin("profile", "p", "u.id=p.user_id")
	LeftJoin(tableName string, aliasName string, condition string) QueryBuilder

	// Right join query table to the select query
	//
	// Ex.: qb.Select("users").Alias("u").RightJoin("profile", "p", "u.id=p.user_id")
	RightJoin(tableName string, aliasName string, condition string) QueryBuilder

	// Order results of select query
	//
	// Ex.: qb.Select("users").OrderBy("name", "ASC")
	OrderBy(column string, direction string) QueryBuilder

	// Grouping results of select query
	//
	// Ex.: qb.Select("users").Columns("id", "MAX(account) AS max_acc").GroupBy("id").Having("max_acc")
	GroupBy(args ...string) QueryBuilder

	// Having clause of select query
	//
	// Ex.: qb.Select("users").Columns("id", "MAX(account) AS max_acc").GroupBy("id").Having("max_acc")
	Having(having string, args ...any) QueryBuilder

	// Limit clause for select query
	//
	// Ex.: qb.Select("users").Limit(10).Offset(5)
	Limit(limit int) QueryBuilder

	// Offset clause for select query
	//
	// Ex.: qb.Select("users").Limit(10).Offset(5)
	Offset(offset int) QueryBuilder

	// Executes query and returns first row
	//
	// Ex.: qb.Select("users").Where("id=?", 5).Row()
	Row() (jvm.M, error)

	// Executes query and returns slice of rows map column/value
	//
	// Ex.: qb.Select("users").Where("id > ?", 5).Rows()
	Rows() ([]jvm.M, error)

	// Executes insert query and returns inserted row identificator with name colID
	//
	// Ex.: qb.Insert("users").Columns("name", "email").Parameters("jv", "jv19841202@gmail.com").ExecReturnID()
	ExecReturnID(colID string) (interface{}, error)

	// Executes insert or update query
	//
	// Ex.: qb.Update("users").Columns("name", "email").Parameters("jv", "jv19841202@gmail.com").Where("name=?", "jv").Exec()
	Exec() error
}

// Constuctor for simple query builder
func New(db *sql.DB) QueryBuilder {
	return &builder{
		db:          db,
		isManualSQL: false,
		params:      make([]any, 0),
		columns:     make([]string, 0),
		joins:       make([]string, 0),
		orderBy:     make([]string, 0),
	}
}

func (b *builder) Alias(alias string) QueryBuilder {
	b.tableAlias = alias
	return b
}

func (b *builder) Select(table string) QueryBuilder {
	b.action = selectAction
	b.tableName = table
	return b
}
func (b *builder) Insert(table string) QueryBuilder {
	b.action = insertAction
	b.tableName = table
	return b
}
func (b *builder) Update(table string) QueryBuilder {
	b.action = updateAction
	b.tableName = table
	return b
}
func (b *builder) Delete(table string) QueryBuilder {
	b.action = deleteAction
	b.tableName = table
	return b
}
func (b *builder) Columns(columns ...string) QueryBuilder {
	b.columns = columns
	return b
}
func (b *builder) Parameters(params ...any) QueryBuilder {
	b.params = params
	return b
}
func (b *builder) ColsWithParams(cols jvm.M) QueryBuilder {
	for key, value := range cols {
		b.columns = append(b.columns, key)
		b.params = append(b.params, value)
	}
	return b
}
func (b *builder) Where(where string, args ...any) QueryBuilder {
	b.where = where
	for _, value := range args {
		b.where = strings.Replace(b.where, "?", "$"+strconv.Itoa(len(b.params)+1), 1)
		b.params = append(b.params, value)
	}
	return b
}
func (b *builder) Having(having string, args ...any) QueryBuilder {
	b.having = "HAVING " + having
	for _, value := range args {
		b.having = strings.Replace(b.having, "?", "$"+strconv.Itoa(len(b.params)+1), 1)
		b.params = append(b.params, value)
	}
	return b
}
func (b *builder) Join(join string, tableName string, aliasName string, condition string) QueryBuilder {
	b.joins = append(b.joins, join+" "+tableName+" AS "+aliasName+" ON "+condition)
	return b
}
func (b *builder) InnerJoin(tableName string, aliasName string, condition string) QueryBuilder {
	b.joins = append(b.joins, "INNER JOIN "+tableName+" AS "+aliasName+" ON "+condition)
	return b
}
func (b *builder) LeftJoin(tableName string, aliasName string, condition string) QueryBuilder {
	b.joins = append(b.joins, "LEFT JOIN "+tableName+" AS "+aliasName+" ON "+condition)
	return b
}
func (b *builder) RightJoin(tableName string, aliasName string, condition string) QueryBuilder {
	b.joins = append(b.joins, "RIGHT JOIN "+tableName+" AS "+aliasName+" ON "+condition)
	return b
}
func (b *builder) OrderBy(column string, direction string) QueryBuilder {
	b.orderBy = append(b.orderBy, "ORDER BY "+column+" "+direction)
	return b
}
func (b *builder) GroupBy(args ...string) QueryBuilder {
	b.groupBy = "GROUP BY " + strings.Join(args, ", ")
	return b
}
func (b *builder) Limit(limit int) QueryBuilder {
	b.limit = limit
	return b
}
func (b *builder) Offset(offset int) QueryBuilder {
	b.offset = offset
	return b
}

func (b *builder) SQL(sql string, args ...any) QueryBuilder {
	b.sql = sql
	b.isManualSQL = true
	for _, value := range args {
		b.sql = strings.Replace(b.sql, "?", "$"+strconv.Itoa(len(b.params)+1), 1)
		b.params = append(b.params, value)
	}
	return b
}

func (b *builder) Row() (jvm.M, error) {
	if !b.isManualSQL {
		if err := b.buildQuery(); err != nil {
			return nil, err
		}
	}
	rows, err := b.db.Query(b.sql, b.params...)
	if err != nil {
		return nil, err
	}
	result := make(jvm.M)
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		row := make([]interface{}, len(columns))
		for i := range row {
			row[i] = new(interface{})
		}
		if err := rows.Scan(row...); err != nil {
			return nil, err
		}

		for i, col := range columns {
			result[col] = *(row[i]).(*interface{})
		}
		return result, nil
	}
	return nil, jve.ErrNotFound
}

func (b *builder) Rows() ([]jvm.M, error) {
	if !b.isManualSQL {
		if err := b.buildQuery(); err != nil {
			return nil, err
		}
	}
	rows, err := b.db.Query(b.sql, b.params...)
	if err != nil {
		return nil, err
	}
	result := make([]jvm.M, 0)
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		row := make([]interface{}, len(columns))
		for i := range row {
			row[i] = new(interface{})
		}
		if err := rows.Scan(row...); err != nil {
			return nil, err
		}
		rows := make(jvm.M)
		for i, col := range columns {
			rows[col] = *(row[i]).(*interface{})
		}
		result = append(result, rows)
	}
	return result, nil
}

func (b *builder) ExecReturnID(colID string) (interface{}, error) {
	if !b.isManualSQL {
		if err := b.buildQuery(); err != nil {
			return nil, err
		}
	}
	b.sql += "\nRETURNING " + colID
	lastInsertedID := new(interface{})
	err := b.db.QueryRow(b.sql, b.params...).Scan(lastInsertedID)
	return lastInsertedID, err
}

func (b *builder) Exec() error {
	if !b.isManualSQL {
		if err := b.buildQuery(); err != nil {
			return err
		}
	}
	_, err := b.db.Exec(b.sql, b.params...)
	return err
}

func (b *builder) buildQuery() error {
	switch b.action {
	case selectAction:
		return b.buildSelect()
	case updateAction:
		return b.buildUpdate()
	case insertAction:
		return b.buildInsert()
	case deleteAction:
		return b.buildDelete()
	}
	return jve.ErrUnknownAction
}

func (b *builder) buildSelect() error {
	b.sql = selectAction + "\n"
	if len(b.columns) > 0 {
		b.sql += " " + strings.Join(b.columns, ", ")
	} else {
		if b.tableAlias == "" {
			b.sql += "*"
		} else {
			b.sql += b.tableAlias + ".*"
		}
	}
	b.sql += "\nFROM " + b.tableName
	if b.tableAlias != "" {
		b.sql += " AS " + b.tableAlias
	}

	if len(b.joins) > 0 {
		for _, j := range b.joins {
			b.sql += "\n" + j + "\n"
		}
	}

	if b.where != "" {
		b.sql += "\nWHERE " + b.where
	}

	if b.groupBy != "" {
		b.sql += "\n" + b.groupBy
	}

	if b.having != "" {
		b.sql += "\n" + b.having
	}

	if b.offset > 0 {
		b.sql += "\nOFFSET " + strconv.Itoa(b.offset)
	}
	if b.limit > 0 {
		b.sql += "\nLIMIT " + strconv.Itoa(b.limit)
	}

	return nil
}

func (b *builder) buildInsert() error {
	// insert into table
	b.sql = insertAction + " INTO " + b.tableName
	if b.tableAlias != "" {
		b.sql += " AS " + b.tableAlias
	}
	b.sql += "\n"

	// (col1, col2)
	if b.tableAlias == "" {
		b.sql += "(" + strings.Join(b.columns, ", ") + ")"
	} else {
		cols := make([]string, len(b.columns))
		for i, c := range b.columns {
			cols[i] = b.tableAlias + "." + c
		}
		b.sql += "(" + strings.Join(cols, ", ") + ")"
	}

	// VALUES
	b.sql += "\nVALUES\n"

	// ($1, $2)
	params := make([]string, len(b.params))
	for i := range b.params {
		params[i] = "$" + strconv.Itoa(i+1)
	}
	b.sql += "(" + strings.Join(params, ", ") + ")"

	return nil
}

func (b *builder) buildUpdate() error {
	// update table as t
	b.sql = "UPDATE " + b.tableName
	if b.tableAlias != "" {
		b.sql += " AS " + b.tableAlias
	}
	// set col1=$1
	b.sql += "\nSET\n"
	if len(b.columns) > len(b.params) {
		return jve.ErrTooManyArgs
	}
	sets := make([]string, len(b.columns))
	for i, col := range b.columns {
		if b.tableAlias == "" {
			sets[i] = col + "=$" + strconv.Itoa(i+1)
		} else {
			sets[i] = b.tableAlias + "." + col + "=$" + strconv.Itoa(i+1)
		}
	}
	b.sql += strings.Join(sets, ",\n")
	if b.where != "" {
		b.sql += "\nWHERE " + b.where
	}

	return nil
}

func (b *builder) buildDelete() error {
	b.sql = "DELETE FROM " + b.tableName
	if b.where != "" {
		b.sql += "\nWHERE " + b.where
	}
	return nil
}
