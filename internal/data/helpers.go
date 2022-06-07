package data

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrEditConflict      = errors.New("edit conflict")
	ErrRecordNotInserted = errors.New("record not inserted")
)

func calculateOffset(pagesize, page int) int {
	return (page - 1) * pagesize
}
