package cache

import (
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestBadgerClient(t *testing.T) {
	key := "testkey"
	value := "testvalue"

	t.Run("ReadValue", func(t *testing.T) {
		testClient, err := NewInMemoryBadgerClient()
		assert.Nil(t, err)

		err = testClient.Set(key, value)
		assert.Nil(t, err)

		sut, err := testClient.Read(key)
		assert.Equal(t, value, string(sut))
		testClient.Close()
	})

	t.Run("ReadEmptyValue", func(t *testing.T) {
		testClient, err := NewInMemoryBadgerClient()
		assert.Nil(t, err)

		sut, err := testClient.Read(key)
		assert.NotNil(t, err)
		assert.Nil(t, sut)

		testClient.Close()
	})

	t.Run("AbilityToOpenTwoDatabases", func(t *testing.T) {
		uuid := uuid2.New().String()
		badgerPath := filepath.Join("/tmp", "badger", uuid)
		testClient, err := NewBadgerClient(badgerPath)
		assert.Nil(t, err)

		sut, err := testClient.Read(key)
		assert.NotNil(t, err)
		assert.Nil(t, sut)

		newUUID := uuid2.New().String()
		newBadgerPath := filepath.Join("/tmp", "badger", newUUID)
		newTestClient, err := NewBadgerClient(newBadgerPath)
		newSut, err := newTestClient.Read(key)
		assert.NotNil(t, err)
		assert.Nil(t, newSut)

		testClient.Close()
		newTestClient.Close()
	})
}
