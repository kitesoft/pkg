// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dml

import (
	"database/sql/driver"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type driverValueBytes []byte

// Value implements the driver.Valuer interface.
func (a driverValueBytes) Value() (driver.Value, error) {
	return []byte(a), nil
}

type driverValueNotSupported uint8

// Value implements the driver.Valuer interface.
func (a driverValueNotSupported) Value() (driver.Value, error) {
	return uint8(a), nil
}

type driverValueNil uint8

// Value implements the driver.Valuer interface.
func (a driverValueNil) Value() (driver.Value, error) {
	return nil, nil
}

type driverValueError uint8

// Value implements the driver.Valuer interface.
func (a driverValueError) Value() (driver.Value, error) {
	return nil, errors.Aborted.Newf("WE've aborted something")
}

func TestArguments_Length_and_Stringer(t *testing.T) {
	t.Parallel()

	t.Run("no slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint(9898).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(MakeNullString("eCom3")).NullInt64(MakeNullInt64(4)).NullFloat64(MakeNullFloat64(2.7)).
			NullBool(MakeNullBool(true)).NullTime(MakeNullTime(now()))
		assert.Exactly(t, 15, args.Len(), "Length mismatch")

		// like fmt.GoStringer
		assert.Exactly(t,
			"dml.MakeArgs(15).Null().Int(-1).Int64(1).Uint64(9898).Uint64(2).Float64(3.100000).Bool(true).String(\"eCom1\").Bytes([]byte{0x65, 0x43, 0x6f, 0x6d, 0x32}).Time(time.Unix(1136228645,2)).NullString(dml.MakeNullString(`eCom3`)).NullInt64(dml.MakeNullInt64(4)).NullFloat64(dml.MakeNullFloat64(2.7)).NullBool(dml.MakeNullBool(true)).NullTime(dml.MakeNullTime(time.Unix(1136228645,2))",
			fmt.Sprintf("%#v", args))
	})

	t.Run("no slices, nulls invalid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64(1).Uint64(2).Float64(3.1).Bool(true).String("eCom1").Bytes([]byte(`eCom2`)).Time(now()).
			NullString(MakeNullString("eCom3", false)).NullInt64(MakeNullInt64(4, false)).NullFloat64(MakeNullFloat64(2.7, false)).
			NullBool(MakeNullBool(true, false)).NullTime(MakeNullTime(now(), false))
		assert.Exactly(t, 14, args.Len(), "Length mismatch")
		assert.Exactly(t,
			"dml.MakeArgs(14).Null().Int(-1).Int64(1).Uint64(2).Float64(3.100000).Bool(true).String(\"eCom1\").Bytes([]byte{0x65, 0x43, 0x6f, 0x6d, 0x32}).Time(time.Unix(1136228645,2)).NullString(dml.NullString{}).NullInt64(dml.NullInt64{}).NullFloat64(dml.NullFloat64{}).NullBool(dml.NullBool{}).NullTime(dml.NullTime{})",
			fmt.Sprintf("%#v", args))
	})

	t.Run("slices, nulls valid", func(t *testing.T) {
		args := MakeArgs(10).
			Null().Int(-1).Int64s(1, 2).Uints(567, 765).Uint64s(2).Float64s(1.2, 3.1).Bools(false, true).Strings("eCom1", "eCom11").BytesSlice(nil, []byte(`eCom2`)).Times(now(), now()).
			NullStrings(MakeNullString("eCom3"), MakeNullString("eCom3")).NullInt64s(MakeNullInt64(4), MakeNullInt64(4)).NullFloat64s(MakeNullFloat64(2.7), MakeNullFloat64(2.7)).
			NullBools(MakeNullBool(true)).NullTimes(MakeNullTime(now()), MakeNullTime(now()))
		assert.Exactly(t, 26, args.Len(), "Length mismatch")
		assert.Exactly(t,
			"dml.MakeArgs(15).Null().Int(-1).Int64s([]int64{1, 2}...).Uints([]uint{0x237, 0x2fd}...).Uint64s([]uint64{0x2}...).Float64s([]float64{1.2, 3.1}...).Bools([]bool{false, true}...).Strings(\"eCom1\",\"eCom11\").BytesSlice([]byte(nil),[]byte{0x65, 0x43, 0x6f, 0x6d, 0x32}).Times(time.Unix(1136228645,2),time.Unix(1136228645,2)).NullStrings(dml.MakeNullString(`eCom3`),dml.MakeNullString(`eCom3`)).NullInt64s(dml.MakeNullInt64(4),dml.MakeNullInt64(4)).NullFloat64s(dml.MakeNullFloat64(2.7),dml.MakeNullFloat64(2.7)).NullBools(dml.MakeNullBool(true)).NullTimes(dml.MakeNullTime(time.Unix(1136228645,2),dml.MakeNullTime(time.Unix(1136228645,2))",
			fmt.Sprintf("%#v", args))
	})
}

func TestIFaceToArgs(t *testing.T) {
	t.Parallel()
	t.Run("not supported", func(t *testing.T) {
		_, err := iFaceToArgs(arguments{}, time.Minute)
		assert.True(t, errors.Is(err, errors.NotSupported), "err should have kind errors.NotSupported")
	})
	t.Run("all types", func(t *testing.T) {
		nt := now()
		args, err := iFaceToArgs(arguments{},
			float32(2.3), float64(2.2),
			int64(5), int(6), int32(7), int16(8), int8(9),
			uint32(math.MaxUint32), uint16(math.MaxUint16), uint8(math.MaxUint8),
			true, "Gopher", []byte(`Hello`),
			now(), &nt, nil,
		)
		require.NoError(t, err)

		assert.Exactly(t, []interface{}{
			float64(2.299999952316284), float64(2.2),
			int64(5), int64(6), int64(7), int64(8), int64(9),
			int64(math.MaxUint32), int64(math.MaxUint16), int64(math.MaxUint8),
			true, "Gopher", []uint8{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			now(), now(), nil,
		}, args.Interfaces())
	})
}

func TestArguments_Clone(t *testing.T) {
	t.Parallel()

	args := MakeArgs(2).Int64(1).String("S1").arguments
	args2 := args.Clone()
	args2[0].value = int(1)
	args2[1].value = "S1a"

	assert.Exactly(t, "dml.MakeArgs(2).Int64(1).String(\"S1\")", args.GoString())
	assert.Exactly(t, "dml.MakeArgs(2).Int(1).String(\"S1a\")", args2.GoString())
}
