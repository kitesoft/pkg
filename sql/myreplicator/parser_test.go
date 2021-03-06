package myreplicator

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestNewBinlogParser_IndexOutOfRange(t *testing.T) {
	parser := NewBinlogParser()

	parser.format = &FormatDescriptionEvent{
		Version:                0x4,
		ServerVersion:          []uint8{0x35, 0x2e, 0x36, 0x2e, 0x32, 0x30, 0x2d, 0x6c, 0x6f, 0x67, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		CreateTimestamp:        0x0,
		EventHeaderLength:      0x13,
		EventTypeHeaderLengths: []uint8{0x38, 0xd, 0x0, 0x8, 0x0, 0x12, 0x0, 0x4, 0x4, 0x4, 0x4, 0x12, 0x0, 0x0, 0x5c, 0x0, 0x4, 0x1a, 0x8, 0x0, 0x0, 0x0, 0x8, 0x8, 0x8, 0x2, 0x0, 0x0, 0x0, 0xa, 0xa, 0xa, 0x19, 0x19, 0x0},
		ChecksumAlgorithm:      0x1,
	}

	parser.tables = map[uint64]*TableMapEvent{
		0x3043b:  {tableIDSize: 6, TableID: 0x3043b, Flags: 0x1, Schema: []uint8{0x73, 0x65, 0x69, 0x75, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72}, Table: []uint8{0x61, 0x70, 0x70, 0x5f, 0x63, 0x72, 0x6f, 0x6e}, ColumnCount: 0x15, ColumnType: []uint8{0x3, 0xf, 0xc, 0xc, 0xf, 0x3, 0xc, 0x3, 0xfc, 0xf, 0x1, 0xfe, 0x2, 0xc, 0xf, 0xf, 0xc, 0xf, 0xf, 0x3, 0xf}, ColumnMeta: []uint16{0x0, 0x180, 0x0, 0x0, 0x2fd, 0x0, 0x0, 0x0, 0x2, 0x180, 0x0, 0xfe78, 0x0, 0x0, 0x180, 0x180, 0x0, 0x180, 0x180, 0x0, 0x2fd}, NullBitmap: []uint8{0xf8, 0xfb, 0x17}},
		0x30453:  {tableIDSize: 6, TableID: 0x30453, Flags: 0x1, Schema: []uint8{0x73, 0x65, 0x69, 0x75, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72}, Table: []uint8{0x73, 0x74, 0x67, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x75, 0x70}, ColumnCount: 0x36, ColumnType: []uint8{0x3, 0x3, 0x3, 0x3, 0x3, 0xf, 0xf, 0x8, 0x3, 0x3, 0x3, 0xf, 0xf, 0x1, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xfe, 0x12, 0xf, 0xf, 0xf, 0xf6, 0x1, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xfe, 0xf6, 0x12, 0x3, 0xf, 0xf, 0x1, 0x1, 0x12, 0xf, 0xf, 0xf, 0xf, 0x3, 0xf, 0x3}, ColumnMeta: []uint16{0x0, 0x0, 0x0, 0x0, 0x0, 0x2fd, 0x12c, 0x0, 0x0, 0x0, 0x0, 0x180, 0x180, 0x0, 0x30, 0x180, 0x180, 0x180, 0x30, 0xc0, 0xfe03, 0x0, 0x180, 0x180, 0x180, 0xc02, 0x0, 0x5a, 0x5a, 0x5a, 0x5a, 0x2fd, 0x2fd, 0x2fd, 0xc0, 0x12c, 0x30, 0xc, 0xfe06, 0xb02, 0x0, 0x0, 0x180, 0x180, 0x0, 0x0, 0x0, 0x180, 0x180, 0x2d, 0x2fd, 0x0, 0x2fd, 0x0}, NullBitmap: []uint8{0xee, 0xdf, 0xff, 0xff, 0xff, 0xff, 0x17}},
		0x30504:  {tableIDSize: 6, TableID: 0x30504, Flags: 0x1, Schema: []uint8{0x73, 0x65, 0x69, 0x75, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72}, Table: []uint8{0x6c, 0x6f, 0x67, 0x5f, 0x73, 0x74, 0x67, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x75, 0x70}, ColumnCount: 0x13, ColumnType: []uint8{0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0x3, 0xf, 0xc, 0xc, 0xc, 0xf, 0xf, 0x3, 0xf}, ColumnMeta: []uint16{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x180, 0x0, 0x0, 0x0, 0x180, 0x180, 0x0, 0x2fd}, NullBitmap: []uint8{0x6, 0xfb, 0x5}},
		0x30450:  {tableIDSize: 6, TableID: 0x30450, Flags: 0x1, Schema: []uint8{0x73, 0x65, 0x69, 0x75, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72}, Table: []uint8{0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74}, ColumnCount: 0x16, ColumnType: []uint8{0x3, 0xfc, 0xc, 0x3, 0xc, 0xf, 0x3, 0xf, 0xc, 0xf, 0xf, 0xf, 0xf, 0x3, 0xc, 0xf, 0xf, 0xf, 0xf, 0x3, 0x3, 0xf}, ColumnMeta: []uint16{0x0, 0x2, 0x0, 0x0, 0x0, 0x2d, 0x0, 0x180, 0x0, 0x180, 0x180, 0x2fd, 0x2d, 0x0, 0x0, 0x180, 0x180, 0x2fd, 0x2d, 0x0, 0x0, 0x2fd}, NullBitmap: []uint8{0xfe, 0xff, 0x2f}},
		0x305bb:  {tableIDSize: 6, TableID: 0x305bb, Flags: 0x1, Schema: []uint8{0x79, 0x6d, 0x63, 0x61, 0x63, 0x68, 0x67, 0x6f}, Table: []uint8{0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x6c, 0x6f, 0x67}, ColumnCount: 0x11, ColumnType: []uint8{0x3, 0x3, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xc, 0xf, 0xf, 0xc, 0xf, 0xf, 0x3, 0xf}, ColumnMeta: []uint16{0x0, 0x0, 0x2fd, 0x12c, 0x2fd, 0x2fd, 0x2d, 0x12c, 0x2fd, 0x0, 0x180, 0x180, 0x0, 0x180, 0x180, 0x0, 0x2fd}, NullBitmap: []uint8{0xfe, 0x7f, 0x1}},
		0x16c36b: {tableIDSize: 6, TableID: 0x16c36b, Flags: 0x1, Schema: []uint8{0x61, 0x63, 0x70}, Table: []uint8{0x73, 0x74, 0x67, 0x5f, 0x6d, 0x61, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6c, 0x69, 0x63, 0x6b, 0x32}, ColumnCount: 0xe, ColumnType: []uint8{0x8, 0x8, 0x3, 0x3, 0x2, 0x2, 0xf, 0x12, 0xf, 0xf, 0x12, 0xf, 0xf, 0xf}, ColumnMeta: []uint16{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2d, 0x0, 0x180, 0x180, 0x0, 0x180, 0x180, 0x2fd}, NullBitmap: []uint8{0xba, 0x3f}},
		0x16c368: {tableIDSize: 6, TableID: 0x16c368, Flags: 0x1, Schema: []uint8{0x73, 0x65, 0x69, 0x75, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72}, Table: []uint8{0x73, 0x74, 0x67, 0x5f, 0x6d, 0x61, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6c, 0x69, 0x63, 0x6b, 0x32}, ColumnCount: 0xe, ColumnType: []uint8{0x8, 0x8, 0x3, 0x3, 0x2, 0x2, 0xf, 0x12, 0xf, 0xf, 0x12, 0xf, 0xf, 0xf}, ColumnMeta: []uint16{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2d, 0x0, 0x180, 0x180, 0x0, 0x180, 0x180, 0x2fd}, NullBitmap: []uint8{0xba, 0x3f}},
		0x3045a:  {tableIDSize: 6, TableID: 0x3045a, Flags: 0x1, Schema: []uint8{0x73, 0x65, 0x69, 0x75, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72}, Table: []uint8{0x63, 0x6f, 0x6e, 0x73}, ColumnCount: 0x1e, ColumnType: []uint8{0x3, 0x3, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xfe, 0x12, 0xf, 0xf, 0xf, 0xf6, 0xf, 0xf, 0xf, 0xf, 0x1, 0x1, 0x1, 0x12, 0xf, 0xf, 0x12, 0xf, 0xf, 0x3, 0xf, 0x1}, ColumnMeta: []uint16{0x0, 0x0, 0x30, 0x180, 0x180, 0x180, 0x30, 0xc0, 0xfe03, 0x0, 0x180, 0x180, 0x180, 0xc02, 0x180, 0x180, 0x180, 0x180, 0x0, 0x0, 0x0, 0x0, 0x180, 0x180, 0x0, 0x180, 0x180, 0x0, 0x2fd, 0x0}, NullBitmap: []uint8{0xfc, 0xff, 0xe3, 0x37}},
		0x3045f:  {tableIDSize: 6, TableID: 0x3045f, Flags: 0x1, Schema: []uint8{0x73, 0x65, 0x69, 0x75, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72}, Table: []uint8{0x63, 0x6f, 0x6e, 0x73, 0x5f, 0x61, 0x64, 0x64, 0x72}, ColumnCount: 0x19, ColumnType: []uint8{0x3, 0x3, 0x3, 0x1, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xfe, 0x3, 0xc, 0x1, 0xc, 0xf, 0xf, 0xc, 0xf, 0xf, 0x3, 0xf, 0x4, 0x4}, ColumnMeta: []uint16{0x0, 0x0, 0x0, 0x0, 0x2fd, 0x2fd, 0x2fd, 0xc0, 0x12c, 0x30, 0xc, 0xfe06, 0x0, 0x0, 0x0, 0x0, 0x180, 0x180, 0x0, 0x180, 0x180, 0x0, 0x2fd, 0x4, 0x4}, NullBitmap: []uint8{0xf0, 0xef, 0x5f, 0x0}},
		0x3065f:  {tableIDSize: 6, TableID: 0x3065f, Flags: 0x1, Schema: []uint8{0x73, 0x65, 0x69, 0x75, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72}, Table: []uint8{0x63, 0x6f, 0x6e, 0x73, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x70, 0x65, 0x61, 0x6b, 0x6f, 0x75, 0x74, 0x5f, 0x6c, 0x65, 0x74, 0x74, 0x65, 0x72}, ColumnCount: 0xd, ColumnType: []uint8{0x3, 0x3, 0x3, 0x3, 0x1, 0x12, 0xf, 0xf, 0x12, 0xf, 0xf, 0x3, 0xf}, ColumnMeta: []uint16{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x180, 0x180, 0x0, 0x180, 0x180, 0x0, 0x2fd}, NullBitmap: []uint8{0xe0, 0x17}},
	}

	_, err := parser.Parse([]byte{
		/* 0x00, */ 0xc1, 0x86, 0x8e, 0x55, 0x1e, 0xa5, 0x14, 0x80, 0xa, 0x55, 0x0, 0x0, 0x0, 0x7, 0xc,
		0xbf, 0xe, 0x0, 0x0, 0x5f, 0x6, 0x3, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x0, 0xd, 0xff,
		0x0, 0x0, 0x19, 0x63, 0x7, 0x0, 0xca, 0x61, 0x5, 0x0, 0x5e, 0xf7, 0xc, 0x0, 0xf5, 0x7,
		0x0, 0x0, 0x1, 0x99, 0x96, 0x76, 0x74, 0xdd, 0x10, 0x0, 0x73, 0x69, 0x67, 0x6e, 0x75, 0x70,
		0x5f, 0x64, 0x62, 0x5f, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6, 0x0, 0x73, 0x79, 0x73, 0x74,
		0x65, 0x6d, 0xb1, 0x3c, 0x38, 0xcb,
	})

	assert.NoError(t, err, "%+v", err)
}
