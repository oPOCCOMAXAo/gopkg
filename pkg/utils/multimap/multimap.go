package multimap

import "strconv"

// MultiMap contains values of int64, float64, string, bool types.
type MultiMap struct {
	data map[string]any
}

func New() *MultiMap {
	return &MultiMap{
		data: make(map[string]any),
	}
}

func (m *MultiMap) set(key string, value any) {
	if m.data == nil {
		m.data = make(map[string]any)
	}

	m.data[key] = value
}

func (m *MultiMap) get(key string) (any, bool) {
	value, ok := m.data[key]

	return value, ok
}

func (m *MultiMap) GetString(name string) string {
	value, ok := m.get(name)
	if !ok {
		return ""
	}

	switch typedValue := value.(type) {
	case string:
		return typedValue
	case int64:
		return strconv.FormatInt(typedValue, 10)
	case float64:
		return strconv.FormatFloat(typedValue, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typedValue)
	}

	return ""
}

func (m *MultiMap) GetInt(name string) int64 {
	value, ok := m.get(name)
	if !ok {
		return 0
	}

	switch typedValue := value.(type) {
	case int64:
		return typedValue
	case float64:
		return int64(typedValue)
	case string:
		i, err := strconv.ParseInt(typedValue, 10, 64)
		if err == nil {
			return i
		}
	case bool:
		if typedValue {
			return 1
		}

		return 0
	}

	return 0
}

func (m *MultiMap) GetFloat(name string) float64 {
	value, ok := m.get(name)
	if !ok {
		return 0
	}

	switch typedValue := value.(type) {
	case float64:
		return typedValue
	case int64:
		return float64(typedValue)
	case string:
		f, err := strconv.ParseFloat(typedValue, 64)
		if err == nil {
			return f
		}
	case bool:
		if typedValue {
			return 1.0
		}

		return 0.0
	}

	return 0
}

func (m *MultiMap) GetBool(name string) bool {
	value, ok := m.get(name)
	if !ok {
		return false
	}

	switch typedValue := value.(type) {
	case bool:
		return typedValue
	case int64:
		return typedValue != 0
	case float64:
		return typedValue != 0.0
	case string:
		b, err := strconv.ParseBool(typedValue)
		if err == nil {
			return b
		}
	}

	return false
}

func (m *MultiMap) SetAny(key string, value any) *MultiMap {
	m.set(key, value)

	return m
}

func (m *MultiMap) SetString(key string, value string) *MultiMap {
	m.set(key, value)

	return m
}

func (m *MultiMap) SetInt(key string, value int64) *MultiMap {
	m.set(key, value)

	return m
}

func (m *MultiMap) SetFloat(key string, value float64) *MultiMap {
	m.set(key, value)

	return m
}

func (m *MultiMap) SetBool(key string, value bool) *MultiMap {
	m.set(key, value)

	return m
}
