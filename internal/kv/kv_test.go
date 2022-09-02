package kv_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/bernata/kvstore/internal/kv"
	"github.com/stretchr/testify/require"
)

func TestKeyNotFound(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(-1)

	_, ok := store.Get("foo/bar")
	require.False(t, ok)
}

func TestRetrieveKey(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(-1)

	err := store.Write("foo/bar", "v1")
	require.NoError(t, err)
	result, ok := store.Get("foo/bar")
	require.True(t, ok)
	require.Equal(t, "v1", result)
}

func TestOverwriteKey(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(-1)

	err := store.Write("foo/bar", "v1")
	require.NoError(t, err)
	result, ok := store.Get("foo/bar")
	require.True(t, ok)
	require.Equal(t, "v1", result)

	err = store.Write("foo/bar", "v2")
	require.NoError(t, err)
	result, ok = store.Get("foo/bar")
	require.True(t, ok)
	require.Equal(t, "v2", result)
}

func TestMultipleKeys(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(-1)

	err := store.Write("foo/bar1", "v1")
	require.NoError(t, err)
	err = store.Write("foo/bar2", "v2")
	require.NoError(t, err)

	result, ok := store.Get("foo/bar1")
	require.True(t, ok)
	require.Equal(t, "v1", result)

	result, ok = store.Get("foo/bar2")
	require.True(t, ok)
	require.Equal(t, "v2", result)
}

func TestDeleteNotExistKey(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(20)

	store.Delete("foo/bar")
	_, ok := store.Get("foo/bar")
	require.False(t, ok)
}

func TestDeleteKey(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(20)

	err := store.Write("foo/bar", "v1")
	require.NoError(t, err)
	store.Delete("foo/bar")
	_, ok := store.Get("foo/bar")
	require.False(t, ok)
}

func TestConcurrentWrites(t *testing.T) {
	t.Parallel()
	const concurrencyCount = 200
	const key = "key"
	wg := sync.WaitGroup{}
	wg.Add(concurrencyCount)
	store := kv.NewStore(20)

	for i := 0; i < concurrencyCount; i++ {
		go func(n int) {
			defer wg.Done()
			err := store.Write(key, fmt.Sprintf("v%d", n))
			require.NoError(t, err)
		}(i)
	}

	wg.Wait()

	result, ok := store.Get(key)
	require.True(t, ok)
	require.True(t, strings.HasPrefix(result, "v"))
}

func TestKeySize(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(20)

	err := store.Write(strings.Repeat("0", 250), "v1")
	require.NoError(t, err)

	err = store.Write(strings.Repeat("0", 251), "v1")
	require.Error(t, err)
}

func TestValueSize(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(20)

	err := store.Write("key1", strings.Repeat("v", 1024*1024))
	require.NoError(t, err)

	err = store.Write("key2", strings.Repeat("v", 1024*1024+1))
	require.Error(t, err)
}
