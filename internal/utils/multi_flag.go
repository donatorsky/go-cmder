package utils

import (
	"fmt"
)

type setter[T any] func(value string) (T, error)

var StringSetter = func(value string) (string, error) {
	return value, nil
}

func NewUniqueMultiFlag[T any](setter setter[T], options ...uniqueSliceOption) *UniqueMultiFlag[T] {
	return &UniqueMultiFlag[T]{
		UniqueSlice: NewUniqueSlice[T](options...),
		setter:      setter,
	}
}

type UniqueMultiFlag[T any] struct {
	*UniqueSlice[T]

	setter setter[T]
}

func (f *UniqueMultiFlag[T]) String() string {
	return fmt.Sprintf("%#v", f.Items())
}

func (f *UniqueMultiFlag[T]) Set(value string) error {
	result, err := f.setter(value)
	if err != nil {
		return err
	}

	_, err = f.Append(result)

	return err
}
