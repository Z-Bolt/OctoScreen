package dataModels

import (
	// "encoding/json"
	"strconv"
	"strings"
	"time"
)

type JsonTime struct{ time.Time }

func (t JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t.Time).Unix(), 10)), nil
}

func (t *JsonTime) UnmarshalJSON(s []byte) (err error) {
	r := strings.Replace(string(s), `"`, ``, -1)
	if r == "null" {
		return nil
	}

	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}

	t.Time = time.Unix(q, 0)
	return
}
