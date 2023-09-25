package db

import "strings"

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	if strings.Index(err.Error(), "Duplicate") >= 0 {
		return true
	}
	return false
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if strings.Index(err.Error(), "sql: no rows in result set") >= 0 {
		return true
	}
	return false
}

func IsDataTooLongError(err error) bool {
	// eg. create app failed: Error 1406: Data too long for column 'name' at row 1
	if err == nil {
		return false
	}
	if strings.Index(err.Error(), "Data too long") >= 0 {
		return true
	}
	return false
}
