package util

import (
	"fmt"
	"time"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	s := t.Format("2006-01-02")
	return []byte(`"` + s + `"`), nil
}

func (d *Date) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*d = Date(v)
		return nil
	case nil:
		*d = Date(time.Time{})
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Date", value)
	}
}
