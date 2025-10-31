package dht

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"testing"
	"time"
)

var (
	pk []byte
	sk []byte
)

func init() {
	pk, _ = base64.StdEncoding.DecodeString("p+i11iGyUqGI/FC4+Nn11HK/S0ZFQ5jjFbsdD8mZa4E=")
	sk, _ = base64.StdEncoding.DecodeString("IXrwlQ+O74nfgtRzXk3vAI9oa0R1C/oHQQFQnb7QlL+n6LXWIbJSoYj8ULj42fXUcr9LRkVDmOMVux0PyZlrgQ==")
}

func TestSplitValue(t *testing.T) {
	var value = []byte("13CYlWV7Perqk75MiTPSrv3obJvhwKPy")
	var sig = []byte("VNWFkOemsPNOjh2qFjhlLLLrITt90tBS5VNWLIzxD4ZAEEWIdb0645e2dKfuTD0k")

	combined := append(value, sig...)

	v, s := splitValue(combined)

	if !bytes.Equal(v, value) {
		t.Fatalf("wanted %s got %s", value, v)
	}

	if !bytes.Equal(s, sig) {
		t.Fatalf("wanted %s got %s", sig, s)
	}
}

type storeMock struct {
	getKey   []byte
	getValue []byte
	getTTL   time.Duration

	getError error
}

func (s *storeMock) Get(out []byte, key []byte) ([]byte, error) {
	return nil, nil
}

func (s *storeMock) Put(key []byte, value []byte, ttl time.Duration) error {
	s.getKey = key
	s.getValue = value
	s.getTTL = ttl
	return s.getError
}

func TestStoreHappyPath(t *testing.T) {
	sm := storeMock{}
	sut := New(Options{
		Store: &sm,
	})

	sig := ed25519.Sign(sk, []byte("hello"))
	var msg = append([]byte("hello"), sig...)

	err := sut.Store(pk, msg)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(msg, sm.getValue) {
		t.Fatalf("expected store.Get(_, value) value to be %x got %x", msg, sm.getValue)
	}

	if !bytes.Equal(pk, sm.getKey) {
		t.Fatalf("expected store.Get(key, _) value to be %x got %x", pk, sm.getKey)
	}

	if sm.getTTL != expire {
		t.Fatalf("expected ttl to be %s got %s", expire, sm.getTTL)
	}
}

func TestStoreErrors(t *testing.T) {
	sut := New(Options{})

	bigKey := make([]byte, keySize+1)
	smallKey := make([]byte, keySize-1)
	value := make([]byte, 64)
	smallValue := make([]byte, 63)
	bigValue := make([]byte, maxValueSize+1)

	err := sut.Store(smallKey, value)
	if err != ErrKeyWrongSize {
		t.Fatalf("expected error '%s' got '%s'", ErrKeyWrongSize, err)
	}

	err = sut.Store(bigKey, value)
	if err != ErrKeyWrongSize {
		t.Fatalf("expected error '%s' got '%s'", ErrKeyWrongSize, err)
	}

	err = sut.Store(pk, smallValue)
	if err != ErrValueTooSmall {
		t.Fatalf("expected error '%s' got '%s'", ErrValueTooSmall, err)
	}

	err = sut.Store(pk, bigValue)
	if err != ErrValueTooLarge {
		t.Fatalf("expected error '%s' got '%s'", ErrValueTooLarge, err)
	}

	err = sut.Store(pk, value)
	if err != ErrSignatureInvalid {
		t.Fatalf("expected error '%s' got '%s'", ErrSignatureInvalid, err)
	}

}
