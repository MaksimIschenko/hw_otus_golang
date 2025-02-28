package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(2)

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		_, ok := c.Get("a")
		require.False(t, ok)

		val, ok := c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)
	})

	t.Run("evict the least recently used element", func(t *testing.T) {
		c := NewCache(3)

		// Add 3 elements
		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		// Use "a" element to make it the most recently used
		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, 1, val)

		// Add 1 more element
		c.Set("d", 4)

		// Check that "b" element was evicted
		val, ok = c.Get("b")
		require.False(t, ok, "element 'b' should be evicted")
		require.Nil(t, val)

		// Check that "c" element is still in cache
		val, ok = c.Get("c")
		require.True(t, ok, "element 'c' should be in cache")
		require.Equal(t, 3, val)

		// Check that "a" element is still in cache
		val, ok = c.Get("a")
		require.True(t, ok, "element 'a' should be in cache")
		require.Equal(t, 1, val)

		// Check that "d" element is in cache
		val, ok = c.Get("d")
		require.True(t, ok, "element 'd' should be in cache")
		require.Equal(t, 4, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
