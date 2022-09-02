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

	store.Write("foo/bar", "v1")
	result, ok := store.Get("foo/bar")
	require.True(t, ok)
	require.Equal(t, "v1", result)
}

func TestOverwriteKey(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(-1)

	store.Write("foo/bar", "v1")
	result, ok := store.Get("foo/bar")
	require.True(t, ok)
	require.Equal(t, "v1", result)

	store.Write("foo/bar", "v2")
	result, ok = store.Get("foo/bar")
	require.True(t, ok)
	require.Equal(t, "v2", result)
}

func TestMultipleKeys(t *testing.T) {
	t.Parallel()
	store := kv.NewStore(-1)

	store.Write("foo/bar1", "v1")
	store.Write("foo/bar2", "v2")

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

	store.Write("foo/bar", "v1")
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
			store.Write(key, fmt.Sprintf("v%d", n))
		}(i)
	}

	wg.Wait()

	result, ok := store.Get(key)
	require.True(t, ok)
	require.True(t, strings.HasPrefix(result, "v"))
}
