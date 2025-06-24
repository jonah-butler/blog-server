package handlers

import (
	ck "blog-api/contextkeys"
	br "blog-api/repositories/blog"
	ur "blog-api/repositories/user"

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

/*
ParseBlogQueryParams modifies a pointer to BlogQuery in place.
It looks up the required query values from the url.Values map,
parsing those values and loading them into the BlogQuery struct
if present, otherwise using the default initialized values from
the struct.
*/
func ParseBlogQueryParams(q *br.BlogQuery, v url.Values) {
	offset := v.Get("offset")

	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err == nil {
			q.Offset = parsedOffset
		}
	}
}

func WriteJSONErr(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encoderErr := json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()}); encoderErr != nil {
		handleFallBackResponse(w)
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encoderErr := json.NewEncoder(w).Encode(data); encoderErr != nil {
		handleFallBackResponse(w)
	}
}

func handleFallBackResponse(w http.ResponseWriter) {
	fallback := `{"error": "Internal Server Error - failed to encode response"}`

	if _, writeErr := w.Write([]byte(fallback)); writeErr != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error: failed to write response: %v", writeErr), http.StatusInternalServerError)
		return
	}
}

func EmptyResponse() map[string]interface{} {
	emptyResponse := map[string]interface{}{}
	return emptyResponse
}

func GetAuthorID(ctx context.Context) (string, bool) {
	authorID, ok := ctx.Value(ck.UserIDKey).(string)
	return authorID, ok
}

func ValidateRequestMime(contentType, mimeType string) bool {
	return strings.Contains(contentType, mimeType)
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

func ParseMultiPartFormBlogUpdate(reader *multipart.Reader) (*br.UpdateBlogInput, error) {
	input := &br.UpdateBlogInput{}
	err := ParseMultiPartForm(reader, input)
	return input, err
}

func ParseMultiPartFormBlogCreate(reader *multipart.Reader) (*br.CreateBlogInput, error) {
	input := &br.CreateBlogInput{}
	err := ParseMultiPartForm(reader, input)
	return input, err
}

func ParseMultiPartFormUserUpdate(reader *multipart.Reader) (*ur.UserUpdatePost, error) {
	input := &ur.UserUpdatePost{}
	err := ParseMultiPartForm(reader, input)
	return input, err
}
