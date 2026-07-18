package stack

import "context"

type contextKey struct{}

type Stack []string

func With(ctx context.Context, s Stack) context.Context {
	return context.WithValue(ctx, contextKey{}, s)
}

func Get(ctx context.Context) Stack {
	if s, ok := ctx.Value(contextKey{}).(Stack); ok {
		return s
	}

	return nil
}

// Push appends name to the stack, creating a copy.
// Trims oldest entries beyond maxLen.
func Push(ctx context.Context, name string, maxLen int) context.Context {
	s := Get(ctx)

	start := max(0, min(len(s)+1-maxLen, len(s)))

	newStack := make(Stack, len(s)+1-start)
	copy(newStack, s[start:])
	newStack[len(newStack)-1] = name

	return With(ctx, newStack)
}
