package main

import (
	pokecache "cli-pokedex/internal"
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestAddOverwrite(t *testing.T) {
	const interval = 5 * time.Second
	cache := pokecache.NewCache(interval)
	cache.Add("https://example.com", []byte("testdata"))
	cache.Add("https://example.com", []byte("moretestdata"))
	val, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}
	if string(val) != "moretestdata" {
		t.Errorf("expected to find newer value, got %s", string(val))
		return
	}
}

func TestGetNotFound(t *testing.T) {
	const interval = 5 * time.Second
	cache := pokecache.NewCache(interval)

	_, ok := cache.Get("https://example.com/not-exists")
	if ok {
		t.Errorf("expected to not find key")
	}
}

func TestReapKeepFresh(t *testing.T) {
	const baseTime = 10 * time.Millisecond
	cache := pokecache.NewCache(baseTime)

	cache.Add("old_key", []byte("old_data"))

	time.Sleep(baseTime / 2)
	cache.Add("new_key", []byte("new_data"))

	time.Sleep(baseTime/2 + 2*time.Millisecond)

	if _, ok := cache.Get("old_key"); ok {
		t.Errorf("expected old_key to be reaped")
	}
	if _, ok := cache.Get("new_key"); !ok {
		t.Errorf("expected new_key to still exist")
	}
}

func TestReapConcurrent(t *testing.T) {
	const interval = 5 * time.Minute
	cache := pokecache.NewCache(interval)

	for i := 0; i < 10; i++ {
		go func(n int) {
			key := fmt.Sprintf("https://example.com/%d", n)
			cache.Add(key, []byte("data"))
			cache.Get(key)
		}(i)
	}
}
