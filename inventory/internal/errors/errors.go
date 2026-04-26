package errs

import "errors"

var (

	// Ошибки деталей.
	ErrPartNotFound = errors.New("деталь не найдена")
	ErrOutOfStock   = errors.New("деталь отсутствует на складе")

	// Ошибки валидации.
	ErrInvalidUUID = errors.New("неверный формат UUID")
)
