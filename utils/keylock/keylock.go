package keylock

import (
	"slices"
	"strings"
	"sync"
)

type KeyLock struct {
	mu    sync.Mutex
	locks map[string]chan struct{}
}

func NewKeyLock() *KeyLock {
	return &KeyLock{
		locks: make(map[string]chan struct{}),
	}
}

func (l *KeyLock) getLockChans(keys []string) []chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	chans := make([]chan struct{}, 0, len(keys))
	for _, key := range keys {
		if ch, ok := l.locks[key]; ok {
			chans = append(chans, ch)
		} else {
			l.locks[key] = make(chan struct{}, 1)
			l.locks[key] <- struct{}{}
			chans = append(chans, l.locks[key])
		}
	}

	return chans
}

func (l *KeyLock) getUnlockFun(lockedChans []chan struct{}) func() {
	return func() {
		for i := len(lockedChans) - 1; i >= 0; i-- {
			lockedChans[i] <- struct{}{}
		}
	}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	keysToLock := make([]string, len(keys))
	copy(keysToLock, keys)
	slices.SortFunc(keysToLock, strings.Compare)

	lockChans := l.getLockChans(keysToLock)
	lockedChans := make([]chan struct{}, 0, len(lockChans))

	for _, ch := range lockChans {
		select {
		case <-cancel:
			l.getUnlockFun(lockedChans)()
			return true, func() {}
		case <-ch:
			lockedChans = append(lockedChans, ch)
		}
	}

	return false, l.getUnlockFun(lockedChans)
}
