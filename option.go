package option

import "errors"

type Option[T any] struct {
	data     T
	has_data bool
}

// Some returns an Option with some value of type T
func Some[T any](data T) Option[T] {
	return Option[T]{
		data:     data,
		has_data: true,
	}
}

// None returns an Option with no value
func None[T any]() Option[T] {
	var t T
	return Option[T]{
		data:     t,
		has_data: false,
	}
}

// IsSome returns `true` if the option is a `Some` value
func (o Option[T]) IsSome() bool {
	return o.has_data
}

// IsSomeAnd returns `true` if the option is a `Some` value
// and the value inside of it matches a predicate
func (o Option[T]) IsSomeAnd(f func(T) bool) bool {
	return o.has_data && f(o.data)
}

// IsNone returns `true` if the option is a `None` value
func (o Option[T]) IsNone() bool {
	return !o.has_data
}

// Expect returns the contained `Some` value.
// Returns an error with the given message if `None`
func (o Option[T]) Expect(msg string) (T, error) {
	if o.IsNone() {
		var t T
		return t, errors.New(msg)
	}
	return o.data, nil
}

// Unwrap returns the contained `Some` value.
// Returns an error if the option is `None`
func (o Option[T]) Unwrap() (T, error) {
	if o.IsNone() {
		var t T
		return t, errors.New("unwrap called on empty Option")
	}
	return o.data, nil
}

// UnwrapOr returns the contained `Some` value or
// a provided default
func (o Option[T]) UnwrapOr(fallback T) T {
	if o.IsNone() {
		return fallback
	}
	return o.data
}

// UnwrapOrElse returns the contained `Some` value or
// computes it from a closure
func (o Option[T]) UnwrapOrElse(fallbackFn func() T) T {
	if o.IsNone() {
		return fallbackFn()
	}
	return o.data
}

// UnwrapOrDefault returns the contained `Some` value or
// the zero value of type T
func (o Option[T]) UnwrapOrDefault() T {
	if o.IsNone() {
		var t T
		return t
	}
	return o.data
}

// Map maps an Option[T] to an Option[U] by applying a function
// to the contained value if it exists
func Map[T any, U any](o Option[T], f func(T) U) Option[U] {
	if o.IsNone() {
		return None[U]()
	}
	return Some(f(o.data))
}

// Inspect calls the provided closure with the contained value
// if it exists.
func (o Option[T]) Inspect(f func(T)) {
	if o.IsSome() {
		f(o.data)
	}
}

// MapOr returns the provided default result (if `None`),
// or applies a function to the contained value (if `Some`).
func MapOr[T any, U any](o Option[T], fallback U, f func(T) U) U {
	if o.IsNone() {
		return fallback
	}
	return f(o.data)
}

// MapOrElse computes a default function result (if `None`),
// or applies a different function to the contained value (if `Some`).
func MapOrElse[T any, U any](o Option[T], fallbackFn func() U, f func(T) U) U {
	if o.IsNone() {
		return fallbackFn()
	}
	return f(o.data)
}

// And returns None if the option is None, otherwise returns optb.
func And[T any, U any](o Option[T], optB Option[U]) Option[U] {
	if o.IsNone() {
		return None[U]()
	}
	return optB
}

// AndThen returns None if the option is None,
// otherwise calls f with the wrapped value and returns the result.
// Also known as flatmap or monadic bind.
func AndThen[T any, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.IsNone() {
		return None[U]()
	}
	return f(o.data)
}

// Filter returns None if the option is None,
// otherwise calls predicate with the wrapped value and returns:
// - Some(t) if predicate returns true (where t is the wrapped value), and
// - None if predicate returns false.
func (o Option[T]) Filter(f func(T) bool) Option[T] {
	if o.IsNone() || !f(o.data) {
		return None[T]()
	}
	return o
}

// Or returns the option if it contains a value,
// otherwise returns optB.
func (o Option[T]) Or(optB Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return optB
}

// OrElse returns the option if it contains a value,
// otherwise calls `f` and returns the result.
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return f()
}

// Xor returns `Some` if exactly one of self, optB is `Some`,
// otherwise returns `None`.
func (o Option[T]) Xor(optB Option[T]) Option[T] {
	if o.IsSome() && optB.IsNone() {
		return o
	}
	if o.IsNone() && optB.IsSome() {
		return optB
	}
	return None[T]()
}

// Insert inserts value into the option, then returns
// a mutable reference to it. If the option already contains
// a value, the old value is dropped.
func (o *Option[T]) Insert(value T) *T {
	o.data = value
	o.has_data = true
	return &o.data
}

// GetOrInsert inserts value into the option if it is `None`,
// then returns a mutable reference to the contained value.
func (o *Option[T]) GetOrInsert(value T) *T {
	if o.IsNone() {
		o.data = value
		o.has_data = true
	}
	return &o.data
}

// GetOrInsert inserts the zero value of type T into
// the option if it is `None`, then returns a mutable
// reference to the contained value.
func (o *Option[T]) GetOrInsertDefault() *T {
	if o.IsNone() {
		var t T
		o.data = t
		o.has_data = true
	}
	return &o.data
}

// GetOrInsertWith inserts a value computed from `f` into
// the option if it is `None`, then returns a mutable
// reference to the contained value.
func (o *Option[T]) GetOrInsertWith(f func() T) *T {
	if o.IsNone() {
		o.data = f()
		o.has_data = true
	}
	return &o.data
}

func (o *Option[T]) Take() Option[T] {
	taken := o.Copy()
	o.has_data = false
	return taken
}

func (o *Option[T]) Replace(t T) Option[T] {
	old := o.Copy()
	o.data = t
	o.has_data = true
	return old
}

func Contains[T comparable](o Option[T], x T) bool {
	if o.IsNone() {
		return false
	}
	return o.data == x
}

func (o Option[T]) Copy() Option[T] {
	return o
}

func Flatten[T any](o Option[Option[T]]) Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	return o.data
}
