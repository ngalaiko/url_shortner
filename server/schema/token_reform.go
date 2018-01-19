// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package schema

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type userTokenTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *userTokenTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("user_token").
func (v *userTokenTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *userTokenTableType) Columns() []string {
	return []string{"id", "token", "user_id", "expired_at"}
}

// NewStruct makes a new struct for that view or table.
func (v *userTokenTableType) NewStruct() reform.Struct {
	return new(UserToken)
}

// NewRecord makes a new record for that table.
func (v *userTokenTableType) NewRecord() reform.Record {
	return new(UserToken)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *userTokenTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// UserTokenTable represents user_token view or table in SQL database.
var UserTokenTable = &userTokenTableType{
	s: parse.StructInfo{Type: "UserToken", SQLSchema: "", SQLName: "user_token", Fields: []parse.FieldInfo{{Name: "ID", Type: "uint64", Column: "id"}, {Name: "Token", Type: "string", Column: "token"}, {Name: "UserID", Type: "uint64", Column: "user_id"}, {Name: "ExpiredAt", Type: "time.Time", Column: "expired_at"}}, PKFieldIndex: 0},
	z: new(UserToken).Values(),
}

// String returns a string representation of this struct or record.
func (s UserToken) String() string {
	res := make([]string, 4)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "Token: " + reform.Inspect(s.Token, true)
	res[2] = "UserID: " + reform.Inspect(s.UserID, true)
	res[3] = "ExpiredAt: " + reform.Inspect(s.ExpiredAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *UserToken) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.Token,
		s.UserID,
		s.ExpiredAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *UserToken) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.Token,
		&s.UserID,
		&s.ExpiredAt,
	}
}

// View returns View object for that struct.
func (s *UserToken) View() reform.View {
	return UserTokenTable
}

// Table returns Table object for that record.
func (s *UserToken) Table() reform.Table {
	return UserTokenTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *UserToken) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *UserToken) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *UserToken) HasPK() bool {
	return s.ID != UserTokenTable.z[UserTokenTable.s.PKFieldIndex]
}

// SetPK sets record primary key.
func (s *UserToken) SetPK(pk interface{}) {
	if i64, ok := pk.(int64); ok {
		s.ID = uint64(i64)
	} else {
		s.ID = pk.(uint64)
	}
}

// check interfaces
var (
	_ reform.View   = UserTokenTable
	_ reform.Struct = (*UserToken)(nil)
	_ reform.Table  = UserTokenTable
	_ reform.Record = (*UserToken)(nil)
	_ fmt.Stringer  = (*UserToken)(nil)
)

func init() {
	parse.AssertUpToDate(&UserTokenTable.s, new(UserToken))
}