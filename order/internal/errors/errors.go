package errs

import "errors"

var (
	// Ошибки заказов.
	ErrOrderNotFound     = errors.New("заказ не найден")
	ErrOrderItemNotFound = errors.New("позиция заказа не найдена")
	ErrOrderAlreadyPaid  = errors.New("заказ уже оплачен")
	ErrOrderCancelled    = errors.New("заказ отменён")

	// Ошибки деталей.
	ErrPartNotFound      = errors.New("деталь не найдена")
	ErrOutOfStock        = errors.New("деталь отсутствует на складе")
	ErrIncompatibleParts = errors.New("детали несовместимы")

	// Ошибки валидации.
	ErrInvalidUUID          = errors.New("неверный формат UUID")
	ErrInvalidPaymentMethod = errors.New("неверный метод оплаты")
)
