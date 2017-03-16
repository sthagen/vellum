//  Copyright (c) 2017 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vellum

import (
	"bytes"
	"reflect"
	"testing"
)

func TestIterator(t *testing.T) {
	var buf bytes.Buffer
	b, err := New(&buf, nil)
	if err != nil {
		t.Fatalf("error creating builder: %v", err)
	}

	err = insertStringMap(b, smallSample)
	if err != nil {
		t.Fatalf("error building: %v", err)
	}

	err = b.Close()
	if err != nil {
		t.Fatalf("error closing: %v", err)
	}

	fst, err := Load(buf.Bytes())
	if err != nil {
		t.Fatalf("error loading set: %v", err)
	}

	got := map[string]uint64{}
	itr, err := fst.Iterator(nil, nil)
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(smallSample, got) {
		t.Errorf("expected %v, got: %v", smallSample, got)
	}
}

func TestIteratorStartKey(t *testing.T) {
	var buf bytes.Buffer
	b, err := New(&buf, nil)
	if err != nil {
		t.Fatalf("error creating builder: %v", err)
	}

	err = insertStringMap(b, smallSample)
	if err != nil {
		t.Fatalf("error building: %v", err)
	}

	err = b.Close()
	if err != nil {
		t.Fatalf("error closing: %v", err)
	}

	fst, err := Load(buf.Bytes())
	if err != nil {
		t.Fatalf("error loading set: %v", err)
	}

	// with start key < "mon", we should still get it
	got := map[string]uint64{}
	itr, err := fst.Iterator([]byte("a"), nil)
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(smallSample, got) {
		t.Errorf("expected %v, got: %v", smallSample, got)
	}

	// with start key = "mon", we should still get it
	got = map[string]uint64{}
	itr, err = fst.Iterator([]byte("mon"), nil)
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(smallSample, got) {
		t.Errorf("expected %v, got: %v", smallSample, got)
	}

	// with start key > "mon", we don't expect to get it
	expect := map[string]uint64{
		"tues":  smallSample["tues"],
		"thurs": smallSample["thurs"],
		"tye":   smallSample["tye"],
	}
	got = map[string]uint64{}
	itr, err = fst.Iterator([]byte("mona"), nil)
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}

	// with start key > "mon", we don't expect to get it
	expect = map[string]uint64{
		"tues":  smallSample["tues"],
		"thurs": smallSample["thurs"],
		"tye":   smallSample["tye"],
	}
	got = map[string]uint64{}
	itr, err = fst.Iterator([]byte("my"), nil)
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}
}

func TestIteratorEndKey(t *testing.T) {
	var buf bytes.Buffer
	b, err := New(&buf, nil)
	if err != nil {
		t.Fatalf("error creating builder: %v", err)
	}

	err = insertStringMap(b, smallSample)
	if err != nil {
		t.Fatalf("error building: %v", err)
	}

	err = b.Close()
	if err != nil {
		t.Fatalf("error closing: %v", err)
	}

	fst, err := Load(buf.Bytes())
	if err != nil {
		t.Fatalf("error loading set: %v", err)
	}

	// with end key > "tye", we should still get it
	got := map[string]uint64{}
	itr, err := fst.Iterator(nil, []byte("zeus"))
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(smallSample, got) {
		t.Errorf("expected %v, got: %v", smallSample, got)
	}

	// with end key = "tye", we should NOT get it (end key exclusive)
	expect := map[string]uint64{
		"mon":   smallSample["mon"],
		"tues":  smallSample["tues"],
		"thurs": smallSample["thurs"],
	}
	got = map[string]uint64{}
	itr, err = fst.Iterator(nil, []byte("tye"))
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}

	// with start key < "tye", we don't expect to get it
	got = map[string]uint64{}
	itr, err = fst.Iterator(nil, []byte("tv"))
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}
}

func TestIteratorSeek(t *testing.T) {
	var buf bytes.Buffer
	b, err := New(&buf, nil)
	if err != nil {
		t.Fatalf("error creating builder: %v", err)
	}

	err = insertStringMap(b, smallSample)
	if err != nil {
		t.Fatalf("error building: %v", err)
	}

	err = b.Close()
	if err != nil {
		t.Fatalf("error closing: %v", err)
	}

	fst, err := Load(buf.Bytes())
	if err != nil {
		t.Fatalf("error loading set: %v", err)
	}

	// seek past thurs (exactly to tues)
	expect := map[string]uint64{
		"mon":  smallSample["mon"],
		"tues": smallSample["tues"],
		"tye":  smallSample["tye"],
	}
	got := map[string]uint64{}
	itr, err := fst.Iterator(nil, nil)
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val

		if string(key) == "mon" {
			err = itr.Seek([]byte("tue"))
		} else {
			err = itr.Next()
		}
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}

	// similar but seek to something after thurs before tues
	got = map[string]uint64{}
	itr, err = fst.Iterator(nil, nil)
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val

		if string(key) == "mon" {
			err = itr.Seek([]byte("thv"))
		} else {
			err = itr.Next()
		}
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}

	// similar but seek to thurs+suffix
	got = map[string]uint64{}
	itr, err = fst.Iterator(nil, nil)
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val

		if string(key) == "mon" {
			err = itr.Seek([]byte("thursday"))
		} else {
			err = itr.Next()
		}
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}

	// seek past last key (still inside iterator boundaries)
	expect = map[string]uint64{
		"mon": smallSample["mon"],
	}
	got = map[string]uint64{}
	itr, err = fst.Iterator(nil, nil)
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val

		if string(key) == "mon" {
			err = itr.Seek([]byte("zzz"))
		} else {
			err = itr.Next()
		}
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}
}

func TestIteratorSeekOutsideBoundaries(t *testing.T) {
	var buf bytes.Buffer
	b, err := New(&buf, nil)
	if err != nil {
		t.Fatalf("error creating builder: %v", err)
	}

	err = insertStringMap(b, smallSample)
	if err != nil {
		t.Fatalf("error building: %v", err)
	}

	err = b.Close()
	if err != nil {
		t.Fatalf("error closing: %v", err)
	}

	fst, err := Load(buf.Bytes())
	if err != nil {
		t.Fatalf("error loading set: %v", err)
	}

	// first test with boundaries should just see thurs/tues
	expect := map[string]uint64{
		"thurs": smallSample["thurs"],
		"tues":  smallSample["tues"],
	}
	got := map[string]uint64{}
	itr, err := fst.Iterator([]byte("th"), []byte("tuesd"))
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}

	// this time try to seek before the start,
	// still shouldn't see mon
	got = map[string]uint64{}
	itr, err = fst.Iterator([]byte("th"), []byte("tuesd"))
	if err != nil {
		t.Fatalf("error before seeking: %v", err)
	}
	err = itr.Seek([]byte("cat"))
	for err == nil {
		key, val := itr.Current()
		got[string(key)] = val
		err = itr.Next()
	}
	if err != ErrIteratorDone {
		t.Errorf("iterator error: %v", err)
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("expected %v, got: %v", expect, got)
	}

	// this time try to seek past the end
	// should see nothing

	itr, err = fst.Iterator([]byte("th"), []byte("tuesd"))
	if err != nil {
		t.Fatalf("error before seeking: %v", err)
	}
	err = itr.Seek([]byte("ty"))
	if err != ErrIteratorDone {
		t.Fatalf("expected ErrIteratorDone, got %v", err)
	}
}

var key []byte
var val uint64

func BenchmarkFSTIteratorAllInMem(b *testing.B) {
	// first build the FST once
	dataset := thousandTestWords
	randomThousandVals := randomValues(dataset)
	var buf bytes.Buffer
	builder, err := New(&buf, nil)
	if err != nil {
		b.Fatalf("error creating builder: %v", err)
	}
	err = insertStrings(builder, dataset, randomThousandVals)
	if err != nil {
		b.Fatalf("error inserting thousand words: %v", err)
	}
	err = builder.Close()
	if err != nil {
		b.Fatalf("error closing builder: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		fst, err := Load(buf.Bytes())
		if err != nil {
			b.Fatalf("error loading FST: %v", err)
		}

		itr, err := fst.Iterator(nil, nil)
		for err == nil {
			key, val = itr.Current()
			err = itr.Next()
		}
		if err != ErrIteratorDone {
			b.Fatalf("iterator error: %v", err)
		}

		err = fst.Close()
		if err != nil {
			b.Fatalf("error closing FST: %v", err)
		}

	}
}
