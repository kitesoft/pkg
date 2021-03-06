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

// +build csall

package objcache_test

import (
	"context"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/allegro/bigcache"
	"github.com/corestoreio/pkg/storage/objcache"
	"github.com/ugorji/go/codec"
)

func benchmark_country_enc(iterationsSetGet int, opts ...objcache.Option) func(b *testing.B) {
	return func(b *testing.B) {
		p, err := objcache.NewManager(opts...)
		if err != nil {
			b.Fatal(err)
		}
		defer func() {
			if err := p.Close(); err != nil {
				b.Fatal(err)
			}
		}()
		cntry := getTestCountry(b) // type already gob.Registered ...
		const wantCountryISO = "US"
		ctx := context.TODO()
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var i int64
			for pb.Next() {
				key := strconv.FormatInt(i, 10) // 1 alloc
				i++

				if err := p.Set(ctx, key, cntry, nil); err != nil {
					b.Fatalf("%+v", err)
				}
				// Double execution might detect storing of type information in streaming encoders
				for j := 0; j < iterationsSetGet; j++ {
					var newCntry = new(Country)
					if err := p.Get(ctx, key, newCntry, nil); err != nil {
						b.Fatalf("%+v", err)
					}
					if newCntry.Country.IsoCode != wantCountryISO {
						b.Fatalf("Country ISO Code must be %q, Have %q", wantCountryISO, newCntry.Country.IsoCode)
					}
				}
			}
		})
	}
}

func benchmark_stores_enc(iterationsSetGet int, opts ...objcache.Option) func(b *testing.B) {
	return func(b *testing.B) {
		p, err := objcache.NewManager(opts...)
		if err != nil {
			b.Fatal(err)
		}
		defer func() {
			if err := p.Close(); err != nil {
				b.Fatal(err)
			}
		}()
		ts := getTestStores() // type already gob.Registered ...
		const wantStoreCode = "nz"
		ctx := context.TODO()
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			var i int64
			for pb.Next() {
				key := strconv.FormatInt(i, 10) // 1 alloc
				i++

				if err := p.Set(ctx, key, ts, nil); err != nil {
					b.Fatal(err)
				}

				// Double execution might detect storing of type information in streaming encoders
				for j := 0; j < iterationsSetGet; j++ {
					var newTS TableStoreSlice
					if err := p.Get(ctx, key, &newTS, nil); err != nil {
						b.Fatal(err)
					}
					if have := newTS[5].Code; have != wantStoreCode {
						b.Fatalf("Store Code in slice position 5 must be %q, Have %q", wantStoreCode, have)
					}
				}
			}
		})
	}
}

func Benchmark_BigCache_Country(b *testing.B) {
	b.Run("Gob_1x", benchmark_country_enc(1, objcache.WithBigCache(bigcache.Config{}), objcache.WithPooledEncoder(gobCodec{}, Country{})))
	b.Run("Gob_2x", benchmark_country_enc(2, objcache.WithBigCache(bigcache.Config{}), objcache.WithPooledEncoder(gobCodec{}, Country{})))
	b.Run("JSON_1x", benchmark_country_enc(1, objcache.WithBigCache(bigcache.Config{}), objcache.WithPooledEncoder(JSONCodec{})))
	b.Run("JSON_2x", benchmark_country_enc(2, objcache.WithBigCache(bigcache.Config{}), objcache.WithPooledEncoder(JSONCodec{})))
	b.Run("MsgPack_1x", benchmark_country_enc(1, objcache.WithBigCache(bigcache.Config{}), objcache.WithEncoder(newMsgPackCodec())))
	b.Run("MsgPack_2x", benchmark_country_enc(2, objcache.WithBigCache(bigcache.Config{}), objcache.WithEncoder(newMsgPackCodec())))
}

func Benchmark_BigCache_Stores(b *testing.B) {
	b.Run("Gob_1x", benchmark_stores_enc(1, objcache.WithBigCache(bigcache.Config{}), objcache.WithPooledEncoder(gobCodec{}, TableStoreSlice{})))
	b.Run("Gob_2x", benchmark_stores_enc(2, objcache.WithBigCache(bigcache.Config{}), objcache.WithPooledEncoder(gobCodec{}, TableStoreSlice{})))
	b.Run("JSON_1x", benchmark_stores_enc(1, objcache.WithBigCache(bigcache.Config{}), objcache.WithPooledEncoder(JSONCodec{})))
	b.Run("JSON_2x", benchmark_stores_enc(2, objcache.WithBigCache(bigcache.Config{}), objcache.WithPooledEncoder(JSONCodec{})))
	b.Run("MsgPack_1x", benchmark_stores_enc(1, objcache.WithBigCache(bigcache.Config{}), objcache.WithEncoder(newMsgPackCodec())))
	b.Run("MsgPack_2x", benchmark_stores_enc(2, objcache.WithBigCache(bigcache.Config{}), objcache.WithEncoder(newMsgPackCodec())))
}

func Benchmark_Redis_Gob(b *testing.B) {
	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		b.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}
	b.Run("Country_1x", benchmark_country_enc(1, objcache.WithRedisURL(redConURL), objcache.WithPooledEncoder(gobCodec{}, Country{})))
	b.Run("Country_2x", benchmark_country_enc(2, objcache.WithRedisURL(redConURL), objcache.WithPooledEncoder(gobCodec{}, Country{})))
	b.Run("Stores_1x", benchmark_stores_enc(1, objcache.WithRedisURL(redConURL), objcache.WithPooledEncoder(gobCodec{}, TableStoreSlice{})))
	b.Run("Stores_2x", benchmark_stores_enc(2, objcache.WithRedisURL(redConURL), objcache.WithPooledEncoder(gobCodec{}, TableStoreSlice{})))
}

func Benchmark_Redis_MsgPack(b *testing.B) {
	redConURL := os.Getenv("CS_REDIS_TEST") // redis://127.0.0.1:6379/3
	if redConURL == "" {
		b.Skip(`Skipping live test because environment CS_REDIS_TEST variable not found.
	export CS_REDIS_TEST="redis://127.0.0.1:6379/3"
		`)
	}
	b.Run("Country_1x", benchmark_country_enc(1, objcache.WithRedisURL(redConURL), objcache.WithEncoder(newMsgPackCodec())))
	b.Run("Country_2x", benchmark_country_enc(2, objcache.WithRedisURL(redConURL), objcache.WithEncoder(newMsgPackCodec())))
	b.Run("Stores_1x", benchmark_stores_enc(1, objcache.WithRedisURL(redConURL), objcache.WithEncoder(newMsgPackCodec())))
	b.Run("Stores_2x", benchmark_stores_enc(2, objcache.WithRedisURL(redConURL), objcache.WithEncoder(newMsgPackCodec())))
}

var ugmsgPackHandle codec.MsgpackHandle

// msgPackCodec cannot be pooled because then it uses too much allocs and slows down.
type msgPackCodec struct{}

func newMsgPackCodec() msgPackCodec {
	return msgPackCodec{}
}

// NewEncoder returns a new json encoder which writes to w
func (c msgPackCodec) NewEncoder(w io.Writer) objcache.Encoder {
	return codec.NewEncoder(w, &ugmsgPackHandle)
}

// NewDecoder returns a new json decoder which reads from r
func (c msgPackCodec) NewDecoder(r io.Reader) objcache.Decoder {
	return codec.NewDecoder(r, &ugmsgPackHandle)
}
