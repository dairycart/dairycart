package dairytest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/pkg/errors"
)

func interfaceArgIsNotPointerOrNil(i interface{}) error {
	if i == nil {
		return errors.New("unmarshalBody cannot accept nil values")
	}
	isNotPtr := reflect.TypeOf(i).Kind() != reflect.Ptr
	if isNotPtr {
		return errors.New("unmarshalBody can only accept pointers")
	}
	return nil
}

func unmarshalBody(res *http.Response, dest interface{}) error {
	// These paths should only ever be reached in tests, an should never be encountered by an end user.
	if err := interfaceArgIsNotPointerOrNil(dest); err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, &dest)
	if err != nil {
		return err
	}

	return nil
}
