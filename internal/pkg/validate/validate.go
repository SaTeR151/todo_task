package validate

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	v        *validator.Validate
	initOnce sync.Once
)

func get() *validator.Validate {
	initOnce.Do(func() {
		v = validator.New()
	})

	return v
}

func Struct(obj any) error {
	if err := get().Struct(obj); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
