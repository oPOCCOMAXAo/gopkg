package tasks

import "sync"

type NameLock struct {
	mu   sync.Mutex
	sets map[string]bool
}

func (l *NameLock) TryLock(name string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.sets == nil {
		l.sets = make(map[string]bool)
	}

	if l.sets[name] {
		return false
	}

	l.sets[name] = true

	return true
}

func (l *NameLock) Unlock(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.sets[name] = false
}
