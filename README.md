# Option

This package offers an `Option` generic type based on the `Option` type from Rust. All functions and methods in this package are named based on the methods found on the Rust type ([docs here](https://doc.rust-lang.org/stable/std/option/)).

## Purpose

The only exported type of this package is `Option`, however this type should never be instantiated directly, but rather through the two constructors: `Some(t T)` or `None()`. The former represennts an Option with a value `t` of generic type `T`, and the latter represents the absense of a value. Using an `Option` allows for the idea of something being "nullable" without actually doing null checks or error checks, while keeping full type safety, and without using pointers to allow base types to be null.

To work with `Option`s, you will often inject functionality into the option type rather than immediately trying to pull the type out of the `Option`. For example, if you have a function that calls that may or may not return an `int`. Normally this would be implemented by either returning a `(int, error)` tuple where the error represents the lack of an `int`, or a `*int` could be used with `nil` as the return value when no `int` is present. Using these two paradigms if we wanted to double the value we would have to do something like the following:

```go
// Error method
func FirstCalculation() (int, error) {...}

v1, err := FirstCalculation()
if err != nil {
  return 0, err
}

return v1 * 2, nil

// Pointer method
func SomeCalculation() *int {...}

v1 := SomeCalculation()
if v1 == nil {
  return nil
}

return v1 * 2
```

In both of these cases, we have to propogate the means with which we checked this error despite having to manully handle the error ourselves. With an `Option`, we apply the logic into the `Option` itself like so:

```go
func SomeCalculation() Option[int] {...}

val := SomeCalculation()
doubled := option.Map(val, func(x int) int {x * 2})
return doubled // Option[int]
```

If `SomeCalculation` returned a value, we will double that value, otherwise we will automatically propogate the `None` which was returned.

## Usage

`Option` types should be instantiated via the `Some` and `None` constructors only (though the default value for an `Option` is a `None`). Most features of this package are implemented as methods on the `Option` type, but a few that require a second generic type are implemented as functions instead. A few examples of how to use the package follow, but more examples on the same functionality can be found in the [Rust `std::option` docs](https://doc.rust-lang.org/stable/std/option/), albeit written in Rust.

### Safe division

In standard Go, a division function is liable to panic if a divisor of 0 is passed. Even if this is handled in the division function, it would be necessary to work with a returned error every time the division function is used. With an `Option` type, we can encapsulate this uncertainty in the type system itself.

```go
func Divide(dividend, divisor int) option.Option[int] {
  if divisor == 0 {
    return option.None()
  }
  return option.Some(dividend/divisor)
}

Divide(6, 2) // Some(3)
Divide(6, 0) // None
```

### Chaining maps

The following example shows how multiple functions can be chained together to modify the `Option`. If in either the `couldReturnNone` or `unsafeMap` functions a `None` is returned, then `composeAll` will return `None`. Otherwise, all of the maps will be applied one after the other to create a `Some` with the final value.

```go
func couldReturnNone(int) option.option[int] {}
func firstMap(int) int {}
func secondMap(int) string {}
func unsafeMap(string) option.Option[int] {}
func finalMap(int) int {}

func composeAll(x int) option.Option[int] {
	return option.Map(
		option.AndThen(
			option.Map(
				option.Map(
					couldReturnNone(x), firstMap), secondMap), unsafeMap), finalMap)
}

func main() {
	fmt.Println(composeAll(5))
}
```

### Working with contained value

For the most part, functionality should be injected into the `Option` rather than trying to pull the inner value out of it. This can be achieved via `Map`s, `AndThen`s, `Inspect`s, and many more utility functions. However, if eventually you need to try to get the value out of the `Option`, it will be less convenient than in the Rust counterpart of this package as Go does not have pattern matching. Instead we must use the `IsSome`, `IsNone`, and `UnWrap*` methods. Note that the `Unwrap` method has undefined behavior if called on a `None` type, and should always be guarded by an `IsSome` or `IsNone` check. If having an error returned is preferable, the `Expect` method can be used.

```go
func getAnOption() option.Option[int] {}

func main() {
  opt := getAnOption()
  if opt.IsSome() {
    val := opt.Unwrap()
    fmt.Println(val + 2)
  } else {
    fmt.Println("No value exists")
  }
}
```

## Functions vs Methods

Most features of this package are implemented as methods on an `Option`. However, due to the lack of generic methods on generic types in Go, some of the methods from the Rust library had to be implemented as package functions. The affected functions are:

- `Map`
- `MapOr`
- `MapOrElse`
- `And`
- `AndThen`
- `Contains`
- `Flatten`

## Missing methods Rust's `std::option`

There are quite a few methods from the Rust `std::option` type that are not implemented in this package. These methods should be methods relating to Rust specific language features such as getting a mutable reference, pinned value, or result type conversion. If there are any missng methods that make sense for a Go `Option` type, feel free to leave a Github issue detailing them.
