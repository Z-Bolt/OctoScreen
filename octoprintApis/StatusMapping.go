package octoprintApis

import (
	// "errors"
	"fmt"
	// "io"
	// "io/ioutil"
	// "log"
	// "net/http"
	// "net/url"
	// "time"
)


type StatusMapping map[int]string

func (this *StatusMapping) Error(code int) error {
	err, ok := (*this)[code]
	if ok {
		return fmt.Errorf(err)
	}

	return nil
}
