package blog

import (
	r "blog-api/repositories/blog"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
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

func ParseMultiePartForm(reader *multipart.Reader) (*r.BlogInput, error) {
	input := new(r.BlogInput)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return input, err
		}

		formName := part.FormName()
		if formName == "image" {
			fileBuffer := new(bytes.Buffer)

			size, err := io.Copy(fileBuffer, part)
			if err != nil {
				return input, err
			}

			fileHeader := &multipart.FileHeader{
				Filename: part.FileName(),
				Header:   part.Header,
				Size:     size,
			}

			input.Image = fileHeader

			input.ImageBytes = fileBuffer.Bytes()
			continue
		}

		buf := new(bytes.Buffer)

		_, err = io.Copy(buf, part)
		if err != nil {
			return input, err
		}

		fieldValue := buf.String()

		switch formName {
		case "categories":
			var categories []string

			err := json.Unmarshal([]byte(fieldValue), &categories)
			if err != nil {
				return input, err
			}

			input.Categories = categories
		case "text":
			input.Text = fieldValue

		case "published":
			published, err := strconv.ParseBool(fieldValue)
			if err != nil {
				return input, err
			}

			input.Published = bool(published)
		case "title":
			input.Title = fieldValue
		case "id":
			input.ID = fieldValue
		}
	}

	return input, nil
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
