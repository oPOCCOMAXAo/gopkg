package ginutils

import "github.com/gin-gonic/gin"

// CtxTypedPointer is helper for get-set pointer values to gin.Context.
//
// Example:
//
//	const SomeTypeValue CtxTypedPointer[SomeType] = "some_context_key"
//
//	func SomeHandler(ctx *gin.Context) {
//		// ...
//		SomeTypeValue.Set(ctx, &someValue)
//		// ...
//		someValue := SomeTypeValue.Get(ctx)
//	}
type CtxTypedPointer[T any] string

// NewTyped creates a new typed pointer context key for the given field name.
func NewTyped[T any](field string) CtxTypedPointer[T] {
	return CtxTypedPointer[T](field)
}

// Set stores a pointer value in the context under this key.
func (typed CtxTypedPointer[T]) Set(ctx *gin.Context, value *T) {
	ctx.Set(string(typed), value)
}

// Constant returns a middleware that sets the given value in the context on every request.
func (typed CtxTypedPointer[T]) Constant(value *T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		typed.Set(ctx, value)
	}
}

// GetOK retrieves the pointer value and reports whether it was found and correctly typed.
func (typed CtxTypedPointer[T]) GetOK(ctx *gin.Context) (*T, bool) {
	value, ok := ctx.Get(string(typed))
	if !ok {
		return nil, false
	}

	res, ok := value.(*T)

	return res, ok
}

// Get retrieves the pointer value, or nil if absent.
func (typed CtxTypedPointer[T]) Get(ctx *gin.Context) *T {
	res, _ := typed.GetOK(ctx)

	return res
}

// GetOrZero retrieves the pointer value, or returns a pointer to a zero-valued T if absent.
func (typed CtxTypedPointer[T]) GetOrZero(ctx *gin.Context) *T {
	res, _ := typed.GetOK(ctx)

	if res == nil {
		res = new(T)
	}

	return res
}

// IsEmpty reports whether the key is absent or holds a nil pointer.
func (typed CtxTypedPointer[T]) IsEmpty(ctx *gin.Context) bool {
	res, ok := typed.GetOK(ctx)

	return !ok || res == nil
}

// CtxTypedValue is helper for get-set values to gin.Context.
//
// Example:
//
//	const SomeTypeValue CtxTypedValue[SomeType] = "some_context_key"
//
//	func SomeHandler(ctx *gin.Context) {
//		// ...
//		SomeTypeValue.Set(ctx, someValue)
//		// ...
//		someValue := SomeTypeValue.Get(ctx)
//	}
type CtxTypedValue[T any] string

// NewTypedValue creates a new typed value context key for the given field name.
func NewTypedValue[T any](field string) CtxTypedValue[T] {
	return CtxTypedValue[T](field)
}

// Set stores a value in the context under this key.
func (typed CtxTypedValue[T]) Set(ctx *gin.Context, value T) {
	ctx.Set(string(typed), value)
}

// Constant returns a middleware that sets the given value in the context on every request.
func (typed CtxTypedValue[T]) Constant(value T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		typed.Set(ctx, value)
	}
}

// GetOK retrieves the value and reports whether it was found and correctly typed.
func (typed CtxTypedValue[T]) GetOK(ctx *gin.Context) (T, bool) {
	value, ok := ctx.Get(string(typed))
	if !ok {
		var zero T

		return zero, false
	}

	res, ok := value.(T)

	return res, ok
}

// Get retrieves the value, or the zero value of T if absent.
func (typed CtxTypedValue[T]) Get(ctx *gin.Context) T {
	res, _ := typed.GetOK(ctx)

	return res
}

// IsEmpty reports whether the key is absent from the context.
func (typed CtxTypedValue[T]) IsEmpty(ctx *gin.Context) bool {
	_, ok := typed.GetOK(ctx)

	return !ok
}
