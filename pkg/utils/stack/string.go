package stack

import "context"

type stringCtxKey struct{}

type String []string

func WithString(ctx context.Context, s String) context.Context {
	return context.WithValue(ctx, stringCtxKey{}, s)
}

func GetString(ctx context.Context) String {
	if s, ok := ctx.Value(stringCtxKey{}).(String); ok {
		return s
	}

	return nil
}

// PushString appends name to the stack, creating a copy.
// Trims oldest entries beyond maxLen.
func PushString(ctx context.Context, name string, maxLen int) context.Context {
	s := GetString(ctx)

	start := max(0, min(len(s)+1-maxLen, len(s)))

	newStack := make(String, len(s)+1-start)
	copy(newStack, s[start:])
	newStack[len(newStack)-1] = name

	return WithString(ctx, newStack)
}
