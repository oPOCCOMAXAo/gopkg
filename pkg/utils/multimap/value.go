package multimap

import (
	"context"
)

type Value func(*MultiMap)

func noValue(*MultiMap) {}

func ContextWithValue(
	ctx context.Context,
	values ...Value,
) context.Context {
	if len(values) == 0 {
		return ctx
	}

	store, ok := FromContext(ctx)
	if !ok {
		store = New()
		ctx = ContextWith(ctx, store)
	}

	for _, value := range values {
		value(store)
	}

	return ctx
}

func StringFromContext(ctx context.Context, name string) string {
	store, ok := FromContext(ctx)
	if !ok {
		return ""
	}

	return store.GetString(name)
}

func IntFromContext(ctx context.Context, name string) int64 {
	store, ok := FromContext(ctx)
	if !ok {
		return 0
	}

	return store.GetInt(name)
}

func BoolFromContext(ctx context.Context, name string) bool {
	store, ok := FromContext(ctx)
	if !ok {
		return false
	}

	return store.GetBool(name)
}

func FloatFromContext(ctx context.Context, name string) float64 {
	store, ok := FromContext(ctx)
	if !ok {
		return 0
	}

	return store.GetFloat(name)
}

func StringValue(name string, value string) Value {
	if value == "" {
		return noValue
	}

	return func(store *MultiMap) {
		store.SetString(name, value)
	}
}

func IntValue(name string, value int64) Value {
	if value == 0 {
		return noValue
	}

	return func(store *MultiMap) {
		store.SetInt(name, value)
	}
}

func FloatValue(name string, value float64) Value {
	if value == 0 {
		return noValue
	}

	return func(store *MultiMap) {
		store.SetFloat(name, value)
	}
}

func BoolValue(name string, value bool) Value {
	if !value {
		return noValue
	}

	return func(store *MultiMap) {
		store.SetBool(name, value)
	}
}
