package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_lruCache_Success(t *testing.T) {
	lru := New[string, string](3)

	err := lru.Put("key1", "val1")
	require.NoError(t, err)
	v, err := lru.Get("key1")
	require.NoError(t, err)
	assert.Equal(t, "val1", v)

	err = lru.Put("key2", "val2")
	require.NoError(t, err)
	v, err = lru.Get("key2")
	require.NoError(t, err)
	assert.Equal(t, "val2", v)

	err = lru.Put("key3", "val3")
	require.NoError(t, err)
	v, err = lru.Get("key3")
	require.NoError(t, err)
	assert.Equal(t, "val3", v)

	err = lru.Put("key4", "val4")
	require.NoError(t, err)
	v, err = lru.Get("key4")
	require.NoError(t, err)
	assert.Equal(t, "val4", v)

	v, err = lru.Get("key1")
	assert.ErrorIs(t, err, ErrNotFound)
}
