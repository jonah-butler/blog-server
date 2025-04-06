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

// eventually remove this
func ParseMultiPartFormBlogUpdateOld(reader *multipart.Reader) (*r.UpdateBlogInput, error) {
	input := new(r.UpdateBlogInput)

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
		case "slug":
			input.Slug = fieldValue
		}
	}

	return input, nil
}

// eventually remove this
func ParseMultiPartFormBlogCreateOld(reader *multipart.Reader) (*r.CreateBlogInput, error) {
	input := new(r.CreateBlogInput)

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
		case "slug":
			input.Slug = fieldValue
		}
	}

	return input, nil
}

func ParseMultiPartForm[T any](reader *multipart.Reader, input *T) error {
	value := reflect.ValueOf(input).Elem()

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		formName := part.FormName()

		if formName == "image" {
			fileBuffer := new(bytes.Buffer)
			size, err := io.Copy(fileBuffer, part)
			if err != nil {
				return err
			}

			if size == 0 || part.FileName() == "" {
				continue
			}

			fileHeader := &multipart.FileHeader{
				Filename: part.FileName(),
				Header:   part.Header,
				Size:     size,
			}

			if field := value.FieldByName("Image"); field.IsValid() {
				field.Set(reflect.ValueOf(fileHeader))
			}
			if field := value.FieldByName("ImageBytes"); field.IsValid() {
				field.Set(reflect.ValueOf(fileBuffer.Bytes()))
			}
			continue
		}

		buf := new(bytes.Buffer)

		_, err = io.Copy(buf, part)
		if err != nil {
			return err
		}
		fieldValue := buf.String()

		field := value.FieldByNameFunc(func(name string) bool {
			field, _ := reflect.TypeOf(input).Elem().FieldByName(name)
			return field.Tag.Get("form") == formName
		})

		if !field.IsValid() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(fieldValue)
		case reflect.Bool:
			parsedBool, err := strconv.ParseBool(fieldValue)
			if err != nil {
				return err
			}
			field.SetBool(parsedBool)
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				var parsedSlice []string
				if err := json.Unmarshal([]byte(fieldValue), &parsedSlice); err != nil {
					return err
				}
				field.Set(reflect.ValueOf(parsedSlice))
			}
		}
	}

	return nil
}

func ParseMultiPartFormBlogUpdate(reader *multipart.Reader) (*r.UpdateBlogInput, error) {
	input := &r.UpdateBlogInput{}
	err := ParseMultiPartForm(reader, input)
	return input, err
}

func ParseMultiPartFormBlogCreate(reader *multipart.Reader) (*r.CreateBlogInput, error) {
	input := &r.CreateBlogInput{}
	err := ParseMultiPartForm(reader, input)
	return input, err
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
