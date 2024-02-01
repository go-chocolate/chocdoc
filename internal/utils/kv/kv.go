package kv

import "strings"

type KV interface {
	Get(key string) string
	Gets(key string) []string
	Set(key, value string)
	Add(key string, value ...string)
	Del(key string)
	Has(key string) bool
}

func NewKV(m map[string]string) KV {
	kv := kvMap{}
	for k, v := range m {
		kv.Set(k, v)
	}
	return kv
}

type kvMap map[string][]string

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (kv kvMap) Get(key string) string {
	if kv == nil {
		return ""
	}
	vs := kv[strings.ToLower(key)]
	if len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (kv kvMap) Gets(key string) []string {
	if kv == nil {
		return nil
	}
	return kv[strings.ToLower(key)]
}

// Set sets the key to value. It replaces any existing
// values.
func (kv kvMap) Set(key, value string) {
	kv[strings.ToLower(key)] = []string{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (kv kvMap) Add(key string, values ...string) {
	key = strings.ToLower(key)
	kv[key] = append(kv[key], values...)
}

// Del deletes the values associated with key.
func (kv kvMap) Del(key string) {
	key = strings.ToLower(key)
	delete(kv, key)
}

// Has checks whether a given key is set.
func (kv kvMap) Has(key string) bool {
	key = strings.ToLower(key)
	_, ok := kv[key]
	return ok
}
