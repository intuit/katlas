package util

import (
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"sync"
	"time"
)

// KeyMutex use for lock resource base on specified key
// use sync.Mutex that ensure map can be accessed safely by multi thread
type KeyMutex struct {
	keys        map[interface{}]interface{}
	mutex       *sync.Mutex
	maxDuration time.Duration
	delayFactor int
}

// NewKeyMutex returns KeyMutex instance
func NewKeyMutex(maxDuration time.Duration, delayFactor int) *KeyMutex {
	return &KeyMutex{
		keys:        make(map[interface{}]interface{}),
		mutex:       &sync.Mutex{},
		maxDuration: maxDuration,
		delayFactor: delayFactor,
	}
}

// TryLock try lock by key
func (k KeyMutex) TryLock(key interface{}) bool {
	timeout := time.NewTimer(k.maxDuration)
	for {
		select {
		case <-timeout.C:
			log.Infof("timeout when try acquire lock for %s ", key)
			return false
		default:
			k.mutex.Lock()
			// if key already locked by others, sleep random time then retry
			if _, ok := k.keys[key]; ok {
				k.mutex.Unlock()
				time.Sleep(time.Duration(rand.Intn(k.delayFactor)) * time.Microsecond)
			} else {
				k.keys[key] = struct{}{}
				k.mutex.Unlock()
				return true
			}
		}
	}
}

// Unlock remove key from map
func (k KeyMutex) Unlock(key interface{}) {
	k.mutex.Lock()
	delete(k.keys, key)
	k.mutex.Unlock()
}
