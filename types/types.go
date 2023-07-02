package types

import (
	"bytes"
	"compress/gzip"
	"database/sql/driver"
	"errors"
	"io/ioutil"
)

// GzippedText is a []byte which transparently gzips data being submitted to
// a database and unzips data being Scanned from a database
type GzippedText []byte

func (g GzippedText) Value() (driver.Value, error) {
	b := make([]byte, 0, len(g))
	buf := bytes.NewBuffer(b)
	w := gzip.NewWriter(buf)
	w.Write(g)
	w.Close()

	return buf.Bytes(), nil
}

func (g *GzippedText) Scan(src any) error {
	var source []byte
	switch src := src.(type) {
	case string:
		source = []byte(src)
	case []byte:
		source = src
	default:
		return errors.New("Incompatible type for GzippedText")
	}
	reader, err := gzip.NewReader(bytes.NewReader(source))
	if err != nil {
		return err
	}
	defer reader.Close()
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	*g = GzippedText(b)
	return nil
}
