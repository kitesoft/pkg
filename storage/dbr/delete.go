package dbr

import (
	"database/sql"
	"fmt"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/util/bufferpool"
	"github.com/corestoreio/csfw/util/errors"
)

// DeleteBuilder contains the clauses for a DELETE statement
type DeleteBuilder struct {
	log.Logger
	Execer
	Preparer

	From alias
	WhereFragments
	OrderBys    []string
	LimitCount  uint64
	LimitValid  bool
	OffsetCount uint64
	OffsetValid bool
}

var _ queryBuilder = (*DeleteBuilder)(nil)

// DeleteFrom creates a new DeleteBuilder for the given table
func (sess *Session) DeleteFrom(from ...string) *DeleteBuilder {
	return &DeleteBuilder{
		Logger:         sess.Logger,
		Execer:         sess.cxn.DB,
		Preparer:       sess.cxn.DB,
		From:           NewAlias(from...),
		WhereFragments: make(WhereFragments, 0, 2),
	}
}

// DeleteFrom creates a new DeleteBuilder for the given table
// in the context for a transaction
func (tx *Tx) DeleteFrom(from ...string) *DeleteBuilder {
	return &DeleteBuilder{
		Logger:         tx.Logger,
		Execer:         tx.Tx,
		Preparer:       tx.Tx,
		From:           NewAlias(from...),
		WhereFragments: make(WhereFragments, 0, 2),
	}
}

// Where appends a WHERE clause to the statement whereSqlOrMap can be a
// string or map. If it's a string, args wil replaces any places holders
func (b *DeleteBuilder) Where(args ...ConditionArg) *DeleteBuilder {
	b.WhereFragments = append(b.WhereFragments, newWhereFragments(args...)...)
	return b
}

// OrderBy appends an ORDER BY clause to the statement
func (b *DeleteBuilder) OrderBy(ord string) *DeleteBuilder {
	b.OrderBys = append(b.OrderBys, ord)
	return b
}

// OrderDir appends an ORDER BY clause with a direction to the statement
func (b *DeleteBuilder) OrderDir(ord string, isAsc bool) *DeleteBuilder {
	if isAsc {
		b.OrderBys = append(b.OrderBys, ord+" ASC")
	} else {
		b.OrderBys = append(b.OrderBys, ord+" DESC")
	}
	return b
}

// Limit sets a LIMIT clause for the statement; overrides any existing LIMIT
func (b *DeleteBuilder) Limit(limit uint64) *DeleteBuilder {
	b.LimitCount = limit
	b.LimitValid = true
	return b
}

// Offset sets an OFFSET clause for the statement; overrides any existing OFFSET
func (b *DeleteBuilder) Offset(offset uint64) *DeleteBuilder {
	b.OffsetCount = offset
	b.OffsetValid = true
	return b
}

// ToSql serialized the DeleteBuilder to a SQL string
// It returns the string with placeholders and a slice of query arguments
func (b *DeleteBuilder) ToSql() (string, []interface{}, error) {
	if len(b.From.Expression) == 0 {
		return "", nil, errors.NewEmptyf(errTableMissing)
	}

	var buf = bufferpool.Get()
	defer bufferpool.Put(buf)
	var args []interface{}

	buf.WriteString("DELETE FROM ")
	buf.WriteString(b.From.QuoteAs())

	// Write WHERE clause if we have any fragments
	if len(b.WhereFragments) > 0 {
		buf.WriteString(" WHERE ")
		writeWhereFragmentsToSql(b.WhereFragments, buf, &args)
	}

	// Ordering and limiting
	if len(b.OrderBys) > 0 {
		buf.WriteString(" ORDER BY ")
		for i, s := range b.OrderBys {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(s)
		}
	}

	if b.LimitValid {
		buf.WriteString(" LIMIT ")
		fmt.Fprint(buf, b.LimitCount)
	}

	if b.OffsetValid {
		buf.WriteString(" OFFSET ")
		fmt.Fprint(buf, b.OffsetCount)
	}
	return buf.String(), args, nil
}

// Exec executes the statement represented by the DeleteBuilder
// It returns the raw database/sql Result and an error if there was one
func (b *DeleteBuilder) Exec() (sql.Result, error) {
	sqlStr, args, err := b.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] delete.exec.tosql")
	}

	fullSql, err := Preprocess(sqlStr, args)
	if err != nil {
		return nil, errors.Wrapf(err, "[dbr] delete.exec.interpolate: %q", fullSql)
	}

	if b.Logger.IsInfo() {
		defer log.WhenDone(b.Logger).Info("dbr.DeleteBuilder.Exec.timing", log.String("sql", fullSql))
	}

	result, err := b.Execer.Exec(fullSql)
	if err != nil {
		return result, errors.Wrap(err, "[dbr] delete.exec.Exec")
	}

	return result, nil
}

// Prepare executes the statement represented by the DeleteBuilder. It returns
// the raw database/sql Statement and an error if there was one. Provided
// arguments in the DeleteBuilder are getting ignored.
func (b *DeleteBuilder) Prepare() (*sql.Stmt, error) {
	sqlStr, _, err := b.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] delete.Prepare.tosql")
	}

	if b.Logger.IsInfo() {
		defer log.WhenDone(b.Logger).Info("dbr.DeleteBuilder.Prepare.timing", log.String("sql", sqlStr))
	}

	stmt, err := b.Preparer.Prepare(sqlStr)
	if err != nil {
		return nil, errors.Wrap(err, "[dbr] delete.Prepare.Prepare")
	}

	return stmt, nil
}
