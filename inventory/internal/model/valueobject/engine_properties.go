package valueobject

import (
	"fmt"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
)

// EngineClass — класс двигателя.
type EngineClass string

const (
	// EngineClassA — двигатель класса A.
	EngineClassA EngineClass = "A"
	// EngineClassB — двигатель класса B.
	EngineClassB EngineClass = "B"
	// EngineClassC — двигатель класса C.
	EngineClassC EngineClass = "C"
)

// NewEngineClass создаёт EngineClass с валидацией допустимых значений.
func NewEngineClass(s string) (EngineClass, error) {
	ec := EngineClass(s)

	switch ec {
	case EngineClassA, EngineClassB, EngineClassC:
		return ec, nil
	default:
		return "", fmt.Errorf("недопустимый класс двигателя %q: %w", s, errs.ErrInvalidProperties)
	}
}

// EngineProperties — свойства двигателя (Value Object).
type EngineProperties struct {
	class            EngineClass
	requiredStrength int
}

func (e *EngineProperties) Class() EngineClass    { return e.class }
func (e *EngineProperties) RequiredStrength() int { return e.requiredStrength }

// NewEngineProperties создаёт свойства двигателя. requiredStrength — минимальная прочность корпуса.
func NewEngineProperties(class EngineClass, requiredStrength int) (PartProperties, error) {
	if requiredStrength <= 0 {
		return PartProperties{}, fmt.Errorf("минимальная прочность корпуса должна быть положительной, получено %d: %w", requiredStrength, errs.ErrInvalidProperties)
	}

	return PartProperties{
		engine: &EngineProperties{
			class:            class,
			requiredStrength: requiredStrength,
		},
	}, nil
}
