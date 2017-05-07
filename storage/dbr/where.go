// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbr

import (
	"strings"

	"unicode/utf8"

	"github.com/corestoreio/errors"
)

// Eq is a map Expression -> value pairs which must be matched in a query.
// Joined at AND statements to the WHERE clause. Implements ConditionArg
// interface. Eq = EqualityMap.
type Eq map[string]Argument

func (eq Eq) appendConditions(wfs WhereFragments) WhereFragments {
	for c, arg := range eq {
		if arg == nil {
			arg = ArgNull()
		}
		wfs = append(wfs, &whereFragment{
			Condition: c,
			Arguments: Arguments{arg},
		})
	}
	return wfs
}

// And sets the logical AND operator. Default case.
func (eq Eq) And() ConditionArg {
	return eq
}

// Or not supported
func (eq Eq) Or() ConditionArg {
	return eq
}

const (
	logicalAnd byte = 'a'
	logicalOr  byte = 'o'
	logicalXor byte = 'x'
	logicalNot byte = 'n'
)

type whereFragment struct {
	// Condition can contain either a column name in the form of table.column or
	// just column. Or Condition can contain an expression. Whenever a condition
	// is not a valid identifier we treat it as an expression.
	Condition string
	Arguments Arguments
	Sub       struct {
		// Select adds a sub-select to the where statement. Condition must be either
		// a column name or anything else which can handle the result of a
		// sub-select.
		Select   *Select
		Operator rune
	}
	// Logical states how multiple where statements will be connected.
	// Default to AND. Possible values are a=AND, o=OR, x=XOR, n=NOT
	Logical byte
	Using   []string
}

func (wf *whereFragment) appendConditions(wfs WhereFragments) WhereFragments {
	return append(wfs, wf)
}

// And sets the logical AND operator
func (wf *whereFragment) And() ConditionArg {
	wf.Logical = logicalAnd
	return wf
}

// Or sets the logical OR operator
func (wf *whereFragment) Or() ConditionArg {
	wf.Logical = logicalOr
	return wf
}

// WhereFragments provides a list WHERE resp. ON clauses.
type WhereFragments []*whereFragment

// Conditions iterates over each WHERE fragment and assembles all conditions
// into a new slice.
func (wfs WhereFragments) Conditions() []string {
	c := make([]string, len(wfs))
	for i, w := range wfs {
		c[i] = w.Condition
	}
	return c
}

// ConditionArg used at argument in Where()
type ConditionArg interface {
	appendConditions(WhereFragments) WhereFragments
	And() ConditionArg // And connects next condition via AND
	Or() ConditionArg  // Or connects next condition via OR
}

// Using add syntactic sugar to a JOIN statement: The USING(column_list) clause
// names a list of columns that must exist in both tables. If tables a and b
// both contain columns c1, c2, and c3, the following join compares
// corresponding columns from the two tables:
//	a LEFT JOIN b USING (c1, c2, c3)
// The columns list gets quoted while writing the query string.
func Using(columns ...string) ConditionArg {
	return &whereFragment{
		Using: columns, // gets quoted during writing the query in ToSQL
	}
}

// SubSelect creates a condition for a WHERE or JOIN statement to compare the
// data in `rawStatementOrColumnName` with the returned value/s of the
// sub-select.
func SubSelect(rawStatementOrColumnName string, operator rune, s *Select) ConditionArg {
	wf := &whereFragment{
		Condition: rawStatementOrColumnName,
	}
	wf.Sub.Select = s
	wf.Sub.Operator = operator
	return wf
}

// Condition adds a condition to a WHERE or HAVING statement.
func Condition(rawStatementOrColumnName string, arg ...Argument) ConditionArg {
	return &whereFragment{
		Condition: rawStatementOrColumnName,
		Arguments: arg,
	}
}

// ParenthesisOpen sets an open parenthesis "(". Mostly used for OR conditions
// in combination with AND conditions.
func ParenthesisOpen() ConditionArg {
	return &whereFragment{
		Condition: "(",
	}
}

// ParenthesisClose sets a closing parenthesis ")". Mostly used for OR
// conditions in combination with AND conditions.
func ParenthesisClose() ConditionArg {
	return &whereFragment{
		Condition: ")",
	}
}

func appendConditions(wf WhereFragments, wargs ...ConditionArg) WhereFragments {
	for _, warg := range wargs {
		wf = warg.appendConditions(wf)
	}
	return wf
}

// stmtType enum of j=join, w=where, h=having
func writeWhereFragmentsToSQL(wf WhereFragments, w queryWriter, args Arguments, stmtType byte) (Arguments, error) {
	if len(wf) == 0 {
		return args, nil
	}

	switch stmtType {
	case 'w':
		w.WriteString(" WHERE ")
	case 'h':
		w.WriteString(" HAVING ")
	}

	i := 0
	for _, f := range wf {

		if stmtType == 'j' {
			if len(f.Using) > 0 {
				w.WriteString(" USING (")
				for j, c := range f.Using {
					if j > 0 {
						w.WriteByte(',')
					}
					Quoter.quote(w, c)
				}
				w.WriteByte(')')
				return args, nil // done, only one using allowed
			}
			if i == 0 {
				w.WriteString(" ON ")
			}
		}

		if f.Condition == ")" {
			w.WriteString(f.Condition)
			continue
		}

		if i > 0 {
			// How the WHERE conditions are connected
			switch f.Logical {
			case logicalAnd:
				w.WriteString(" AND ")
			case logicalOr:
				w.WriteString(" OR ")
			case logicalXor:
				w.WriteString(" XOR ")
			case logicalNot:
				w.WriteString(" NOT ")
			default:
				w.WriteString(" AND ")
			}
		}

		if f.Condition == "(" {
			i = 0
			w.WriteString(f.Condition)
			continue
		}

		w.WriteByte('(')
		addArg := false
		if isValidIdentifier(f.Condition) > 0 { // must be an expression
			_, _ = w.WriteString(f.Condition)
			addArg = true
			if len(f.Arguments) == 1 && f.Arguments[0].operator() > 0 {
				writeOperator(w, f.Arguments[0].operator(), true)
			}
		} else {
			Quoter.FquoteAs(w, f.Condition)

			if f.Sub.Select != nil {
				writeOperator(w, f.Sub.Operator, false)
				w.WriteByte('(')
				subArgs, err := f.Sub.Select.toSQL(w)
				w.WriteByte(')')
				if err != nil {
					return nil, errors.Wrapf(err, "[dbr] writeWhereFragmentsToSQL failed SubSelect for table: %q", f.Sub.Select.Table.String())
				}
				args = append(args, subArgs...)
			} else {
				// a column only supports one argument. If not provided we panic
				// with an index out of bounds error.
				addArg = writeOperator(w, f.Arguments[0].operator(), true)
			}
		}
		w.WriteByte(')')

		if addArg {
			args = append(args, f.Arguments...)
		}
		i++
	}
	return args, nil
}

// maxIdentifierLength see http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
const maxIdentifierLength = 64
const dummyQualifier = "X" // just a dummy value, can be optimized later

// IsValidIdentifier checks the permissible syntax for identifiers. Certain
// objects within MySQL, including database, table, index, column, alias, view,
// stored procedure, partition, tablespace, and other object names are known as
// identifiers. ASCII: [0-9,a-z,A-Z$_] (basic Latin letters, digits 0-9, dollar,
// underscore) Max length 63 characters.
//
// Returns 0 if the identifier is valid.
//
// http://dev.mysql.com/doc/refman/5.7/en/identifiers.html
func isValidIdentifier(objectName string) int8 {
	if objectName == sqlStar {
		return 0
	}
	qualifier := dummyQualifier
	if i := strings.IndexByte(objectName, '.'); i >= 0 {
		qualifier = objectName[:i]
		objectName = objectName[i+1:]
	}

	validQualifier := isNameValid(qualifier)
	if validQualifier == 0 && objectName == sqlStar {
		return 0
	}
	if validQualifier > 0 {
		return validQualifier
	}
	return isNameValid(objectName)
}

// isNameValid returns 0 if the name is valid or an error number identifying
// where the name becomes invalid.
func isNameValid(name string) int8 {
	if name == dummyQualifier {
		return 0
	}

	ln := len(name)
	if ln > maxIdentifierLength || name == "" {
		return 1 //errors.NewNotValidf("[csdb] Incorrect identifier. Too long or empty: %q", name)
	}
	pos := 0
	for pos < ln {
		r, w := utf8.DecodeRuneInString(name[pos:])
		pos += w
		if !mapAlNum(r) {
			return 2 // errors.NewNotValidf("[csdb] Invalid character in name %q", name)
		}
	}
	return 0
}

func mapAlNum(r rune) bool {
	var ok bool
	switch {
	case '0' <= r && r <= '9':
		ok = true
	case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z':
		ok = true
	case r == '$', r == '_':
		ok = true
	}
	return ok
}
