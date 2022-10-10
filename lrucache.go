package main

import (
	"container/list"
	"errors"
	"sync"
)

var (
	ErrNotFound          = errors.New("element not found")
	ErrUnexpectedKeyType = errors.New("unexpected key type")
)

type Cache[K comparable, V any] interface {
	Get(key K) (V, error)
	Put(key K, val V) error
	Len() int
}

type (
	lruCache[K comparable, V any] struct {
		size      int
		keysQueue *list.List
		elements  map[K]*node[V]
		lock      *sync.Mutex
	}

	node[V any] struct {
		value  V
		keyPtr *list.Element
	}
)

func New[K comparable, V any](size int) Cache[K, V] {
	return &lruCache[K, V]{
		size:      size,
		keysQueue: list.New(),
		elements:  make(map[K]*node[V], size),
		lock:      &sync.Mutex{},
	}
}

func (lru *lruCache[K, V]) Get(key K) (V, error) {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	elem, ok := lru.elements[key]
	if !ok || elem == nil {
		var res V // default value of generic type V
		return res, ErrNotFound
	}

	lru.keysQueue.MoveToFront(elem.keyPtr)

	return elem.value, nil
}

func (lru *lruCache[K, V]) Put(key K, val V) error {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	if elem, ok := lru.elements[key]; ok && elem != nil {
		elem.value = val
		lru.keysQueue.MoveToFront(elem.keyPtr)
		return nil
	}

	if lru.size == lru.keysQueue.Len() {
		if err := lru.evict(); err != nil {
			return err
		}
	}
	newKeyPtr := lru.keysQueue.PushFront(&key)
	lru.elements[key] = &node[V]{value: val, keyPtr: newKeyPtr}

	return nil
}

func (lru *lruCache[K, V]) Len() int {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	return lru.keysQueue.Len()
}

func (lru *lruCache[K, V]) evict() error {
	keyForEviction := lru.keysQueue.Back()
	key, ok := keyForEviction.Value.(*K)
	if !ok {
		return ErrUnexpectedKeyType
	}
	delete(lru.elements, *key)
	lru.keysQueue.Remove(keyForEviction)
	return nil
}
