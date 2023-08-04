package cdb

import "io"

// Get seeks 'key' in the cdb, similar to Find(), but it resembles the
// interface describe by the cdbget programm
func Get(w io.Writer, c *Cdb, key []byte, skip int) (n int64, err error) {
	var record io.Reader
	c.FindStart()
	for ; skip >= 0; skip-- {
		record, err = c.FindNext(key)
		if err == io.EOF {
			return 0, nil
		} else if err != nil {
			return 0, err
		}
	}
	return io.Copy(w, record)
}
