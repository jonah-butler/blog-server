package blog

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"reflect"
	"strings"
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

func ValidateRequestMime(contentType, mimeType string) bool {
	return strings.Contains(contentType, mimeType)
}

func ParseMultiePartForm(reader *multipart.Reader) error {
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		fmt.Println(part.FormName())
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
