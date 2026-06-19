package domain

import "errors"

// ErrNotFound ошибка в случаи отсутствия данных
// ErrAlreadyExists ошибка в случаи если данные уже существуют
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidPeriod = errors.New("invalid period")
)
