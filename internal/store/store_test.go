package store

import (
	"bytes"
	"reflect"
	"testing"
)

func TestStoreSet(t *testing.T) {
	store := New()

	store.Set("key", []byte("value"))

	got, present := store.Get("key")
	want := []byte("value")
	res := bytes.Compare(got, want)
	if res != 0 {
		t.Errorf("got %s, want %s", got, want)
	}
	if present != true {
		t.Error("got false, want true")
	}
}

func TestStoreGet(t *testing.T) {
	store := New()
	store.Set("key", []byte("value"))

	got, present := store.Get("key")

	want := []byte("value")
	res := bytes.Compare(got, want)
	if res != 0 {
		t.Errorf("got %s, want %s", got, want)
	}
	if present != true {
		t.Error("got false, want true")
	}
}

func TestStoreGetMissingKey(t *testing.T) {
	store := New()

	got, present := store.Get("key")

	want := []byte(nil)
	res := bytes.Compare(got, want)
	if res != 0 {
		t.Errorf("got %s, want %s", got, want)
	}
	if present != false {
		t.Error("got true, want false")
	}
}

func TestStoreDelete(t *testing.T) {
	store := New()
	store.Set("key", []byte("value"))

	store.Delete("key")

	got, present := store.Get("key")
	want := []byte(nil)
	res := bytes.Compare(got, want)
	if res != 0 {
		t.Errorf("got %s, want %s", got, want)
	}
	if present != false {
		t.Error("got true, want false")
	}
}

func TestStoreGetKeys(t *testing.T) {
	store := New()
	store.Set("key2", []byte("value2"))
	store.Set("key1", []byte("value1"))

	got := store.Keys()

	want := []string{"key1", "key2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s, want %s", got, want)
	}
}
