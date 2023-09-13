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
		capacity := 3
		c := NewCache(capacity)
		var val interface{}
		var ok bool
		for i := 0; i < capacity; i++ {
			c.Set(Key("key_"+strconv.Itoa(i)), i)

			val, ok = c.Get(Key("key_" + strconv.Itoa(i)))
			require.True(t, ok)
			require.Equal(t, i, val)
		}

		c.Set(Key("key_"+strconv.Itoa(capacity)), capacity)

		val, ok = c.Get(Key("key_0"))
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get(Key("key_" + strconv.Itoa(capacity)))
		require.True(t, ok)
		require.Equal(t, capacity, val)

		for i := 1; i < capacity; i++ {
			c.Get(Key("key_" + strconv.Itoa(i)))
		}
		c.Set(Key("key_0"), 0)

		val, ok = c.Get(Key("key_0"))
		require.True(t, ok)
		require.Equal(t, 0, val)

		val, ok = c.Get(Key("key_" + strconv.Itoa(capacity)))
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
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

func TestDisplaceLastItem(t *testing.T) {
	c := NewCache(2)
	c.Set(Key("ket_0"), 100)
	c.Set(Key("ket_1"), 200)
	c.(*lruCache).DisplaceLastItem()

	val, ok := c.Get(Key("ket_0"))
	require.False(t, ok)
	require.Nil(t, val)

	val, ok = c.Get(Key("ket_1"))
	require.True(t, ok)
	require.Equal(t, 200, val)
}
