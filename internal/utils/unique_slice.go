package utils

import (
	"fmt"
	"slices"
)

type Uniquer interface {
	UniqueValue() any
}

type uniqueSliceOptions struct {
	length         uint
	capacity       uint
	onDuplicateKey func(key, item any) error
}

type uniqueSliceOption func(options *uniqueSliceOptions)

func NewUniqueSlice[T any](options ...uniqueSliceOption) *UniqueSlice[T] {
	uniqueSliceOptions := &uniqueSliceOptions{
		length:         0,
		capacity:       0,
		onDuplicateKey: func(key, item any) error { return nil },
	}

	for _, option := range options {
		option(uniqueSliceOptions)
	}

	return &UniqueSlice[T]{
		items:   make([]T, uniqueSliceOptions.length, uniqueSliceOptions.capacity),
		cache:   map[any]*T{},
		options: uniqueSliceOptions,
	}
}

func UniqueSliceWithLength(length uint) uniqueSliceOption {
	return func(options *uniqueSliceOptions) {
		options.length = length
	}
}

func UniqueSliceWithCapacity(capacity uint) uniqueSliceOption {
	return func(options *uniqueSliceOptions) {
		options.capacity = capacity
	}
}

func UniqueSliceWithOnDuplicateKeyError(onDuplicateKey func(key, item any) error) uniqueSliceOption {
	return func(options *uniqueSliceOptions) {
		options.onDuplicateKey = onDuplicateKey
	}
}

type UniqueSlice[T any] struct {
	items   []T
	cache   map[any]*T
	options *uniqueSliceOptions
}

func (s *UniqueSlice[T]) Items() []T {
	return s.items
}

func (s *UniqueSlice[T]) Append(item T) (*UniqueSlice[T], error) {
	key, exists := s.has(item)

	if exists {
		return s, s.options.onDuplicateKey(key, item)
	}

	s.items = append(s.items, item)
	s.cache[key] = &item

	return s, nil
}

func (s *UniqueSlice[T]) Has(item T) (exists bool) {
	_, exists = s.has(item)

	return
}

func (s *UniqueSlice[T]) Len() int {
	return len(s.items)
}

func (s *UniqueSlice[T]) Empty() bool {
	return len(s.items) == 0
}

func (s *UniqueSlice[T]) Sort(cmp func(i, j T) int) {
	slices.SortFunc(s.items, cmp)
}

func (s *UniqueSlice[T]) GetByItem(item *T) error {
	key, exist := s.has(*item)
	if !exist {
		return fmt.Errorf("item with key %q does not exist", key)
	}

	*item = *s.cache[key]

	return nil
}

func (s *UniqueSlice[T]) has(item T) (key any, exists bool) {
	key = s.getUniqueKey(item)

	_, exists = s.cache[key]

	return
}

func (s *UniqueSlice[T]) getUniqueKey(v any) any {
	if uniquer, ok := v.(Uniquer); ok {
		return uniquer.UniqueValue()
	}

	return v
}
