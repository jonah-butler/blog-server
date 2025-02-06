package handlers

import (
	"errors"
	"net/url"
	"reflect"
)

// find a better spot for this eventually
func DecodeQueryToStruct(q url.Values, s interface{}) error {

	sEle := reflect.TypeOf(s).Elem()

	sVal := reflect.ValueOf(s)

	if sEle.Kind() != reflect.Struct {
		return errors.New("parameter [s <interface{}>] is not a struct")
	}

	for i := 0; i < sEle.NumField(); i++ {

		param := sEle.Field(i).Tag.Get("param")

		val := q.Get(param)

		if val == "" {

			val = sEle.Field(i).Tag.Get("default")

		}

		sVal.Elem().Field(i).SetString(val)
	}

	return nil
}
