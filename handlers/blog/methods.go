package blog

import (
	"errors"
	"net/url"
	"reflect"
)

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

// func SafeMultiFormParse(reader *multipart.Reader, d interface{}) error {
// 	dType := reflect.TypeOf(d)

// 	if dType.Kind() != reflect.Struct && dType.Kind() != reflect.Ptr {
// 		return fmt.Errorf("not a valid struct")
// 	}

// 	dEle := dType.Elem()
// 	// dVal := reflect.ValueOf(d)

// 	for i := 0; i < dEle.NumField(); i++ {
// 		fmt.Println(dEle.Field(i).Tag.Get("param"))
// 	}

// 	return nil
// }
