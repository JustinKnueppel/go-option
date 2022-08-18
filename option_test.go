package option_test

import (
	"testing"

	"github.com/JustinKnueppel/go-option"
)

func TestIsSome(t *testing.T) {
	tests := map[string]struct {
		value    option.Option[int]
		expected bool
	}{
		"some_value": {
			value:    option.Some(1),
			expected: true,
		},
		"no_value": {
			value:    option.None[int](),
			expected: false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.IsSome() != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestIsSomeAnd(t *testing.T) {
	tests := map[string]struct {
		value     option.Option[int]
		predicate func(int) bool
		expected  bool
	}{
		"some_value_true": {
			value:     option.Some(1),
			predicate: func(x int) bool { return x == 1 },
			expected:  true,
		},
		"some_value_false": {
			value:     option.Some(1),
			predicate: func(x int) bool { return x == 2 },
			expected:  false,
		},
		"no_value": {
			value:     option.None[int](),
			predicate: func(x int) bool { return x == 1 },
			expected:  false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.IsSomeAnd(tc.predicate) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestIsNone(t *testing.T) {
	tests := map[string]struct {
		value    option.Option[int]
		expected bool
	}{
		"some_value": {
			value:    option.Some(1),
			expected: false,
		},
		"no_value": {
			value:    option.None[int](),
			expected: true,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.IsNone() != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestExpect(t *testing.T) {
	tests := map[string]struct {
		value         option.Option[int]
		inner         int
		msg           string
		expectedError bool
	}{
		"some_value": {
			value:         option.Some(1),
			inner:         1,
			msg:           "",
			expectedError: false,
		},
		"no_value": {
			value:         option.None[int](),
			inner:         0,
			msg:           "No value",
			expectedError: true,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			val, err := tc.value.Expect(tc.msg)
			if (err != nil) != tc.expectedError {
				t.Fail()
			}
			if tc.expectedError && tc.msg != err.Error() {
				t.Fail()
			}
			if !tc.expectedError && val != tc.inner {
				t.Fail()
			}
		})
	}
}
func TestUnwrap(t *testing.T) {
	tests := map[string]struct {
		value option.Option[int]
		inner int
	}{
		"some_value": {
			value: option.Some(1),
			inner: 1,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			val := tc.value.Unwrap()
			if val != tc.inner {
				t.Fail()
			}
		})
	}
}
func TestUnwrapOr(t *testing.T) {
	tests := map[string]struct {
		value     option.Option[int]
		alternate int
		inner     int
	}{
		"some_value": {
			value:     option.Some(1),
			alternate: 2,
			inner:     1,
		},
		"no_value": {
			value:     option.None[int](),
			alternate: 2,
			inner:     2,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			val := tc.value.UnwrapOr(tc.alternate)
			if val != tc.inner {
				t.Fail()
			}
		})
	}
}
func TestUnwrapOrElse(t *testing.T) {
	tests := map[string]struct {
		value     option.Option[int]
		alternate func() int
		inner     int
	}{
		"some_value": {
			value:     option.Some(1),
			alternate: func() int { return 2 },
			inner:     1,
		},
		"no_value": {
			value:     option.None[int](),
			alternate: func() int { return 3 },
			inner:     3,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			val := tc.value.UnwrapOrElse(tc.alternate)
			if val != tc.inner {
				t.Fail()
			}
		})
	}
}
func TestUnwrapOrDefault(t *testing.T) {
	tests := map[string]struct {
		value option.Option[int]
		inner int
	}{
		"some_value": {
			value: option.Some(1),
			inner: 1,
		},
		"no_value": {
			value: option.None[int](),
			inner: 0,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			val := tc.value.UnwrapOrDefault()
			if val != tc.inner {
				t.Fail()
			}
		})
	}
}
func TestMapSameType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    option.Option[int]
		function func(int) int
		result   option.Option[int]
	}{
		"some_value": {
			value:    option.Some(2),
			function: func(i int) int { return i * 3 },
			result:   option.Some(6),
		},
		"some_value_closure": {
			value:    option.Some(1),
			function: func(i int) int { return i * closureConstant },
			result:   option.Some(closureConstant),
		},
		"no_value": {
			value:    option.None[int](),
			function: func(i int) int { return i * 3 },
			result:   option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			optB := option.Map(tc.value, tc.function)
			if optB != tc.result {
				t.Fail()
			}
		})
	}
}
func TestMapDifferentType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    option.Option[int]
		function func(int) bool
		result   option.Option[bool]
	}{
		"some_value": {
			value:    option.Some(2),
			function: func(i int) bool { return i == 2 },
			result:   option.Some(true),
		},
		"some_value_closure": {
			value:    option.Some(1),
			function: func(i int) bool { return i > closureConstant },
			result:   option.Some(false),
		},
		"no_value": {
			value:    option.None[int](),
			function: func(i int) bool { return true },
			result:   option.None[bool](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			optB := option.Map(tc.value, tc.function)
			if optB != tc.result {
				t.Fail()
			}
		})
	}
}
func TestInspect(t *testing.T) {
	initialClosureValue := 0
	closureValue := initialClosureValue
	tests := map[string]struct {
		value         option.Option[int]
		function      func(int)
		closureResult int
	}{
		"some_value": {
			value:         option.Some(2),
			function:      func(i int) {},
			closureResult: initialClosureValue,
		},
		"some_value_closure": {
			value:         option.Some(1),
			function:      func(i int) { closureValue = i },
			closureResult: 1,
		},
		"no_value": {
			value:         option.None[int](),
			function:      func(i int) {},
			closureResult: initialClosureValue,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			closureValue = initialClosureValue
			copy := tc.value.Copy()
			val := tc.value.Inspect(tc.function)
			if closureValue != tc.closureResult {
				t.Fail()
			}
			if val != copy {
				t.Fail()
			}
		})
	}
}
func TestMapOrSameType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    option.Option[int]
		fallback int
		function func(int) int
		result   int
	}{
		"some_value": {
			value:    option.Some(2),
			fallback: 1,
			function: func(i int) int { return i + 5 },
			result:   7,
		},
		"some_value_closure": {
			value:    option.Some(2),
			fallback: 1,
			function: func(i int) int { return i * closureConstant },
			result:   10,
		},
		"no_value": {
			value:    option.None[int](),
			fallback: 1,
			function: func(i int) int { return i },
			result:   1,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			optB := option.MapOr(tc.value, tc.fallback, tc.function)
			if optB != tc.result {
				t.Fail()
			}
		})
	}
}

func TestMapOrDifferentType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    option.Option[int]
		fallback bool
		function func(int) bool
		result   bool
	}{
		"some_value": {
			value:    option.Some(2),
			fallback: false,
			function: func(i int) bool { return i < 5 },
			result:   true,
		},
		"some_value_closure": {
			value:    option.Some(2),
			fallback: false,
			function: func(i int) bool { return i > closureConstant },
			result:   false,
		},
		"no_value": {
			value:    option.None[int](),
			fallback: true,
			function: func(i int) bool { return false },
			result:   true,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			optB := option.MapOr(tc.value, tc.fallback, tc.function)
			if optB != tc.result {
				t.Fail()
			}
		})
	}
}
func TestMapOrElseSameType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    option.Option[int]
		fallback func() int
		function func(int) int
		result   int
	}{
		"some_value": {
			value:    option.Some(2),
			fallback: func() int { return 2 },
			function: func(i int) int { return i + 5 },
			result:   7,
		},
		"some_value_closure": {
			value:    option.Some(2),
			fallback: func() int { return 3 },
			function: func(i int) int { return i * closureConstant },
			result:   10,
		},
		"no_value": {
			value:    option.None[int](),
			fallback: func() int { return 3 },
			function: func(i int) int { return i },
			result:   3,
		},
		"no_value_closure": {
			value:    option.None[int](),
			fallback: func() int { return closureConstant },
			function: func(i int) int { return i },
			result:   closureConstant,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			optB := option.MapOrElse(tc.value, tc.fallback, tc.function)
			if optB != tc.result {
				t.Fail()
			}
		})
	}
}

func TestMapOrElseDifferentType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    option.Option[int]
		fallback func() bool
		function func(int) bool
		result   bool
	}{
		"some_value": {
			value:    option.Some(2),
			fallback: func() bool { return true },
			function: func(i int) bool { return i < 5 },
			result:   true,
		},
		"some_value_closure": {
			value:    option.Some(2),
			fallback: func() bool { return closureConstant == 5 },
			function: func(i int) bool { return i > closureConstant },
			result:   false,
		},
		"no_value": {
			value:    option.None[int](),
			fallback: func() bool { return true },
			function: func(i int) bool { return false },
			result:   true,
		},
		"no_value_closure": {
			value:    option.None[int](),
			fallback: func() bool { return closureConstant < 3 },
			function: func(i int) bool { return false },
			result:   false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			optB := option.MapOrElse(tc.value, tc.fallback, tc.function)
			if optB != tc.result {
				t.Fail()
			}
		})
	}
}
func TestAndSameType(t *testing.T) {
	tests := map[string]struct {
		value    option.Option[int]
		other    option.Option[int]
		expected option.Option[int]
	}{
		"some_some": {
			value:    option.Some(1),
			other:    option.Some(2),
			expected: option.Some(2),
		},
		"some_none": {
			value:    option.Some(1),
			other:    option.None[int](),
			expected: option.None[int](),
		},
		"none_some": {
			value:    option.None[int](),
			other:    option.Some(2),
			expected: option.None[int](),
		},
		"none_none": {
			value:    option.None[int](),
			other:    option.None[int](),
			expected: option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if option.And(tc.value, tc.other) != tc.expected {
				t.Fail()
			}
		})
	}
}

func TestAndDifferentType(t *testing.T) {
	tests := map[string]struct {
		value    option.Option[int]
		other    option.Option[string]
		expected option.Option[string]
	}{
		"some_some": {
			value:    option.Some(1),
			other:    option.Some("hello"),
			expected: option.Some("hello"),
		},
		"some_none": {
			value:    option.Some(1),
			other:    option.None[string](),
			expected: option.None[string](),
		},
		"none_some": {
			value:    option.None[int](),
			other:    option.Some("hello"),
			expected: option.None[string](),
		},
		"none_none": {
			value:    option.None[int](),
			other:    option.None[string](),
			expected: option.None[string](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if option.And(tc.value, tc.other) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestAndThenSameType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    option.Option[int]
		function func(int) option.Option[int]
		result   option.Option[int]
	}{
		"some_value_return_some": {
			value:    option.Some(2),
			function: func(i int) option.Option[int] { return option.Some(i * 3) },
			result:   option.Some(6),
		},
		"some_value_return_none": {
			value:    option.Some(2),
			function: func(i int) option.Option[int] { return option.None[int]() },
			result:   option.None[int](),
		},
		"some_value_closure": {
			value:    option.Some(1),
			function: func(i int) option.Option[int] { return option.Some(i * closureConstant) },
			result:   option.Some(closureConstant),
		},
		"no_value_return_some": {
			value:    option.None[int](),
			function: func(i int) option.Option[int] { return option.Some(i * 3) },
			result:   option.None[int](),
		},
		"no_value_return_none": {
			value:    option.None[int](),
			function: func(i int) option.Option[int] { return option.None[int]() },
			result:   option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			optB := option.AndThen(tc.value, tc.function)
			if optB != tc.result {
				t.Fail()
			}
		})
	}
}

func TestAndThenDifferentType(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    option.Option[int]
		function func(int) option.Option[bool]
		result   option.Option[bool]
	}{
		"some_value_return_some": {
			value:    option.Some(2),
			function: func(i int) option.Option[bool] { return option.Some(i == 2) },
			result:   option.Some(true),
		},
		"some_value_return_none": {
			value:    option.Some(2),
			function: func(i int) option.Option[bool] { return option.None[bool]() },
			result:   option.None[bool](),
		},
		"some_value_closure": {
			value:    option.Some(1),
			function: func(i int) option.Option[bool] { return option.Some(i > closureConstant) },
			result:   option.Some(false),
		},
		"no_value_return_some": {
			value:    option.None[int](),
			function: func(i int) option.Option[bool] { return option.Some(true) },
			result:   option.None[bool](),
		},
		"no_value_return_none": {
			value:    option.None[int](),
			function: func(i int) option.Option[bool] { return option.None[bool]() },
			result:   option.None[bool](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			optB := option.AndThen(tc.value, tc.function)
			if optB != tc.result {
				t.Fail()
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := map[string]struct {
		value     option.Option[int]
		predicate func(int) bool
		expected  option.Option[int]
	}{
		"some_value_true": {
			value:     option.Some(1),
			predicate: func(x int) bool { return x == 1 },
			expected:  option.Some(1),
		},
		"some_value_false": {
			value:     option.Some(1),
			predicate: func(x int) bool { return x == 2 },
			expected:  option.None[int](),
		},
		"no_value": {
			value:     option.None[int](),
			predicate: func(x int) bool { return x == 1 },
			expected:  option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.Filter(tc.predicate) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestOr(t *testing.T) {
	tests := map[string]struct {
		value    option.Option[int]
		other    option.Option[int]
		expected option.Option[int]
	}{
		"some_some": {
			value:    option.Some(1),
			other:    option.Some(2),
			expected: option.Some(1),
		},
		"some_none": {
			value:    option.Some(1),
			other:    option.None[int](),
			expected: option.Some(1),
		},
		"none_some": {
			value:    option.None[int](),
			other:    option.Some(2),
			expected: option.Some(2),
		},
		"none_none": {
			value:    option.None[int](),
			other:    option.None[int](),
			expected: option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.Or(tc.other) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestOrElse(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value    option.Option[int]
		other    func() option.Option[int]
		expected option.Option[int]
	}{
		"some_some": {
			value:    option.Some(1),
			other:    func() option.Option[int] { return option.Some(2) },
			expected: option.Some(1),
		},
		"some_none": {
			value:    option.Some(1),
			other:    func() option.Option[int] { return option.None[int]() },
			expected: option.Some(1),
		},
		"none_some": {
			value:    option.None[int](),
			other:    func() option.Option[int] { return option.Some(2) },
			expected: option.Some(2),
		},
		"none_some_closure": {
			value:    option.None[int](),
			other:    func() option.Option[int] { return option.Some(closureConstant) },
			expected: option.Some(closureConstant),
		},
		"none_none": {
			value:    option.None[int](),
			other:    func() option.Option[int] { return option.None[int]() },
			expected: option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.OrElse(tc.other) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestXor(t *testing.T) {
	tests := map[string]struct {
		value    option.Option[int]
		other    option.Option[int]
		expected option.Option[int]
	}{
		"some_some": {
			value:    option.Some(1),
			other:    option.Some(2),
			expected: option.None[int](),
		},
		"some_none": {
			value:    option.Some(1),
			other:    option.None[int](),
			expected: option.Some(1),
		},
		"none_some": {
			value:    option.None[int](),
			other:    option.Some(2),
			expected: option.Some(2),
		},
		"none_none": {
			value:    option.None[int](),
			other:    option.None[int](),
			expected: option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			if tc.value.Xor(tc.other) != tc.expected {
				t.Fail()
			}
		})
	}
}
func TestInsert(t *testing.T) {
	tests := map[string]struct {
		value    option.Option[int]
		newInner int
	}{
		"some_value": {
			value:    option.Some(1),
			newInner: 2,
		},
		"no_value": {
			value:    option.None[int](),
			newInner: 3,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			ref := tc.value.Insert(tc.newInner)
			if *ref != tc.newInner {
				t.Fail()
			}
			if tc.value != option.Some(tc.newInner) {
				t.Fail()
			}
			*ref = *ref + 2
			if tc.value != option.Some(tc.newInner+2) {
				t.Fail()
			}
		})
	}
}
func TestGetOrInsert(t *testing.T) {
	tests := map[string]struct {
		value       option.Option[int]
		newInner    int
		result      option.Option[int]
		resultInner int
	}{
		"some_value": {
			value:       option.Some(1),
			newInner:    2,
			result:      option.Some(1),
			resultInner: 1,
		},
		"no_value": {
			value:       option.None[int](),
			newInner:    3,
			result:      option.Some(3),
			resultInner: 3,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			ref := tc.value.GetOrInsert(tc.newInner)
			if *ref != tc.resultInner {
				t.Fail()
			}
			if tc.value != tc.result {
				t.Fail()
			}
			*ref = *ref + 2
			if tc.value != option.Some(tc.resultInner+2) {
				t.Fail()
			}
		})
	}
}
func TestGetOrInsertDefault(t *testing.T) {
	tests := map[string]struct {
		value       option.Option[int]
		result      option.Option[int]
		resultInner int
	}{
		"some_value": {
			value:       option.Some(1),
			result:      option.Some(1),
			resultInner: 1,
		},
		"no_value": {
			value:       option.None[int](),
			result:      option.Some(0),
			resultInner: 0,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			ref := tc.value.GetOrInsertDefault()
			if *ref != tc.resultInner {
				t.Fail()
			}
			if tc.value != tc.result {
				t.Fail()
			}
			*ref = *ref + 2
			if tc.value != option.Some(tc.resultInner+2) {
				t.Fail()
			}
		})
	}
}
func TestGetOrInsertWith(t *testing.T) {
	closureConstant := 5
	tests := map[string]struct {
		value       option.Option[int]
		fallback    func() int
		result      option.Option[int]
		resultInner int
	}{
		"some_value": {
			value:       option.Some(1),
			fallback:    func() int { return 5 },
			result:      option.Some(1),
			resultInner: 1,
		},
		"no_value": {
			value:       option.None[int](),
			fallback:    func() int { return 3 },
			result:      option.Some(3),
			resultInner: 3,
		},
		"no_value_closure": {
			value:       option.None[int](),
			fallback:    func() int { return closureConstant },
			result:      option.Some(closureConstant),
			resultInner: closureConstant,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			ref := tc.value.GetOrInsertWith(tc.fallback)
			if *ref != tc.resultInner {
				t.Fail()
			}
			if tc.value != tc.result {
				t.Fail()
			}
			*ref = *ref + 2
			if tc.value != option.Some(tc.resultInner+2) {
				t.Fail()
			}
		})
	}
}
func TestTake(t *testing.T) {
	tests := map[string]struct {
		value option.Option[int]
	}{
		"some_value": {
			value: option.Some(1),
		},
		"no_value": {
			value: option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			copy := tc.value.Copy()
			val := tc.value.Take()
			if tc.value != option.None[int]() {
				t.Fail()
			}
			if val != copy {
				t.Fail()
			}
		})
	}
}
func TestReplace(t *testing.T) {
	tests := map[string]struct {
		value    option.Option[int]
		newInner int
		result   option.Option[int]
	}{
		"some_value": {
			value:    option.Some(1),
			newInner: 2,
			result:   option.Some(2),
		},
		"no_value": {
			value:    option.None[int](),
			newInner: 3,
			result:   option.Some(3),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			copy := tc.value.Copy()
			old := tc.value.Replace(tc.newInner)
			if tc.value != tc.result {
				t.Fail()
			}
			if old != copy {
				t.Fail()
			}
		})
	}
}
func TestContains(t *testing.T) {
	tests := map[string]struct {
		value  option.Option[int]
		target int
		result bool
	}{
		"some_value_equals": {
			value:  option.Some(1),
			target: 1,
			result: true,
		},
		"some_value_not_equals": {
			value:  option.Some(1),
			target: 2,
			result: false,
		},
		"no_value": {
			value:  option.None[int](),
			target: 3,
			result: false,
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			contains := option.Contains(tc.value, tc.target)
			if contains != tc.result {
				t.Fail()
			}
		})
	}
}
func TestCopy(t *testing.T) {
	tests := map[string]struct {
		value  option.Option[int]
		result option.Option[int]
	}{
		"some_value": {
			value:  option.Some(1),
			result: option.Some(1),
		},
		"no_value": {
			value:  option.None[int](),
			result: option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			copy := tc.value.Copy()
			if copy != tc.result {
				t.Fail()
			}
			if &copy == &tc.value {
				t.Fail()
			}
		})
	}
}
func TestFlatten(t *testing.T) {
	tests := map[string]struct {
		value  option.Option[option.Option[int]]
		result option.Option[int]
	}{
		"some_some": {
			value:  option.Some(option.Some(1)),
			result: option.Some(1),
		},
		"some_none": {
			value:  option.Some(option.None[int]()),
			result: option.None[int](),
		},
		"none": {
			value:  option.None[option.Option[int]](),
			result: option.None[int](),
		},
	}

	for tname, tc := range tests {
		t.Run(tname, func(t *testing.T) {
			flattened := option.Flatten(tc.value)
			if flattened != tc.result {
				t.Fail()
			}
		})
	}
}
