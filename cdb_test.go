package cdb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

type rec struct {
	key    string
	values []string
}

var records = []rec{
	{"one", []string{"1"}},
	{"two", []string{"2", "22"}},
	{"three", []string{"3", "33", "333"}},
}

var data []byte // set by init()

func init() {
	b := bytes.NewBuffer(nil)
	for _, rec := range records {
		key := rec.key
		for _, value := range rec.values {
			b.WriteString(fmt.Sprintf("+%d,%d:%s->%s\n", len(key), len(value), key, value))
		}
	}
	b.WriteByte('\n')
	data = b.Bytes()
}

func TestCdbMake(t *testing.T) {
	make_from_data(data, t)
}

func TestCdbDump(t *testing.T) {
	f, buf := make_from_data(data, t), bytes.NewBuffer(nil)
	err := Dump(buf, f.bytesReader())
	if err != nil {
		t.Fatalf("Dump failed: %s", err)
	}
	if !bytes.Equal(buf.Bytes(), data) {
		t.Fatalf("Dump round-trip failed")
	}
}

func TestCdbGet(t *testing.T) {
	f := make_from_data(data, t)
	c, buf := New(f.bytesReader()), bytes.NewBuffer(nil)
	for _, rec := range records {
		for skip, val := range rec.values {
			buf.Reset()
			if _, err := Get(buf, c, []byte(rec.key), skip); err != nil {
				t.Fatalf("cdb.Get failed: %s", err)
			}
			if !bytes.Equal(buf.Bytes(), []byte(val)) {
				t.Fatalf("cdb.Get failed: expected %q, got %q", val, buf.Bytes())
			}
			t.Logf("%q => %d: %q", rec.key, skip, val)
		}
	}
}

func TestCdbDataAndFind(t *testing.T) {
	f := make_from_data(data, t)
	c := New(f.bytesReader())

	_, err := c.Data([]byte("does not exist"))
	if err != io.EOF {
		t.Fatalf("non-existent key should return io.EOF")
	}

	for _, rec := range records {
		key := []byte(rec.key)
		values := rec.values

		v, err := c.Data(key)
		if err != nil {
			t.Fatalf("Record read failed: %s", err)
		}

		if !bytes.Equal(v, []byte(values[0])) {
			t.Fatal("Incorrect value returned")
		}
		t.Logf("%q => %q", key, v)

		c.FindStart()
		for _, value := range values {
			sr, err := c.FindNext(key)
			if err != nil {
				t.Fatalf("Record read failed: %s", err)
			}

			data, err := ioutil.ReadAll(sr)
			if err != nil {
				t.Fatalf("Record read failed: %s", err)
			}

			if !bytes.Equal(data, []byte(value)) {
				t.Fatal("value mismatch")
			}

			t.Logf("  %q => %q", key, data)
		}
		// Read all values, so should get EOF
		_, err = c.FindNext(key)
		if err != io.EOF {
			t.Fatalf("Expected EOF, got %s", err)
		}
	}
}

func TestEmptyFile(t *testing.T) {
	f := make_from_data([]byte("\n\n"), t)
	readNum := makeNumReader(f.bytesReader())
	for i := 0; i < 256; i++ {
		_ = readNum() // table pointer
		tableLen := readNum()
		if tableLen != 0 {
			t.Fatalf("table %d has non-zero length: %d", i, tableLen)
		}
	}

	c := New(f.bytesReader())
	_, err := c.Data([]byte("does not exist"))
	if err != io.EOF {
		t.Fatalf("non-existent key should return io.EOF")
	}
}

func make_from_data(d []byte, t *testing.T) *memFile {
	writer := &memFile{}
	if err := Make(writer, bytes.NewBuffer(d)); err != nil {
		t.Fatalf("Make failed: %s", err)
	}
	return writer
}

// 'memFile' is a naive implementation of a io.WriteSeeker (backed by
// a []byte-buffer) to be used in the tests without creating any real files
//
// NOTE: it might be usefull elsewhere, but the .Seek() method might
// move memFile.i (the write-position) behind the len(buf).
// memFile.growIfNeeded takes care of growing the buffer
type memFile struct {
	buf []byte
	i   int64
}

func (f *memFile) Write(data []byte) (int, error) {
	f.growIfNeeded(int64(len(data)))
	n := copy(f.buf[f.i:], data)
	f.i += int64(n)
	return n, nil
}

func (f *memFile) Seek(offset int64, whence int) (abs int64, _ error) {
	switch whence {
	default:
		return 0, errors.New("bufWriteSeeker.Seek: invalid whence")
	case 0:
		abs = offset
	case 1:
		abs = f.i + offset
	case 2:
		abs = int64(len(f.buf)) + offset
	}
	if abs < 0 {
		return 0, errors.New("bufWriteSeeker.Seek: negative position")
	}
	f.i = abs
	return
}

// grows the buffer to hold (mw.i + n) bytes
func (f *memFile) growIfNeeded(n int64) {
	if needed := ((f.i + n) - int64(len(f.buf))); needed > 0 {
		f.buf = append(f.buf, make([]byte, needed)...)
	}
}

func (f *memFile) bytesReader() *bytes.Reader { return bytes.NewReader(f.buf) }
