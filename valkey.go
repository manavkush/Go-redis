package main

import (
	"log/slog"
	"sync"
)

type KV struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewKV() *KV {
	return &KV{
		data: make(map[string][]byte),
	}
}

func (kv *KV) Set(key, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[string(key)] = val
	return nil
}

func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	val, ok := kv.data[string(key)]
	if !ok {
		slog.Error("Key not present in the KV", "key", key)
		return val, ok
	}
	return val, true
}
