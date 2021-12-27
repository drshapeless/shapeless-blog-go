package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readFileAsIntArray(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	temp := []int{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		temp = append(temp, i)
	}

	return temp, nil
}

func isIDUnique(id int, arr []int) bool {
	for _, value := range arr {
		if id == value {
			return false
		}
	}

	return true
}

func isStringExist(s string, arr []string) bool {
	for _, value := range arr {
		if s == value {
			return true
		}
	}
	return false
}

func isStringSlicesEqual(o, n []string) bool {
	if len(o) != len(n) {
		return false
	} else {
		sort.Strings(o)
		sort.Strings(n)

		for i, v1 := range o {
			if n[i] != v1 {
				return false
			}
		}
		return true
	}
}

func removedAndAddedValue(o, n []string) ([]string, []string) {
	// Return two slices of strings, first is removed values, second is
	// added values.

	var same []string

	for _, v1 := range o {
		for _, v2 := range n {
			if v1 == v2 {
				same = append(same, v1)
				break
			}
		}
	}

	return removeStrings(o, same), removeStrings(n, same)
}

func removeStrings(s, r []string) []string {
	ret := s
	for _, v := range r {
		ret = removeString(ret, v)
	}
	return ret
}

func removeString(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func removeInt(arr []int, r int) []int {
	for i, v := range arr {
		if r == v {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
