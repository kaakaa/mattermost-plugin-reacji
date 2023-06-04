package kvstore

import (
	"testing"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"
	"github.com/stretchr/testify/assert"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
)

func TestReacjiStoreGet(t *testing.T) {
	t.Run("all fine", func(t *testing.T) {
		in := &reacji.List{}
		api := &plugintest.API{}
		api.On("KVGet", keyList).Return(in.EncodeToByte(), nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		out, err := store.Reacji().Get()
		assert.NoError(t, err)
		assert.Equal(t, in, out)
	})
	t.Run("there is no existing data", func(t *testing.T) {
		init := &reacji.List{}
		api := &plugintest.API{}
		api.On("KVGet", keyList).Return(nil, nil)
		api.On("KVSet", keyList, init.EncodeToByte()).Return(nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		out, err := store.Reacji().Get()
		assert.NoError(t, err)
		assert.Equal(t, init, out)
	})
	t.Run("KVGet fail", func(t *testing.T) {
		appErr := &model.AppError{}
		api := &plugintest.API{}
		api.On("KVGet", keyList).Return(nil, appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		out, err := store.Reacji().Get()
		assert.Error(t, err)
		assert.Nil(t, out)
	})
	t.Run("KVSet fail", func(t *testing.T) {
		init := &reacji.List{}
		appErr := &model.AppError{}
		api := &plugintest.API{}
		api.On("KVGet", keyList).Return(nil, nil)
		api.On("KVSet", keyList, init.EncodeToByte()).Return(appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		out, err := store.Reacji().Get()
		assert.Error(t, err)
		assert.Nil(t, out)
	})
	t.Run("invalid data", func(t *testing.T) {
		in := []byte{}
		api := &plugintest.API{}
		api.On("KVGet", keyList).Return(in, nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		out, err := store.Reacji().Get()
		assert.Error(t, err)
		assert.Nil(t, out)
	})
}

func TestReacjiStoreUpdate(t *testing.T) {
	t.Run("all fine", func(t *testing.T) {
		prev := &reacji.List{}
		new := &reacji.List{}
		opt := model.PluginKVSetOptions{
			Atomic:   true,
			OldValue: prev.EncodeToByte(),
		}

		api := &plugintest.API{}
		api.On("KVSetWithOptions", keyList, new.EncodeToByte(), opt).Return(true, nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		err := store.Reacji().Update(prev, new)
		assert.NoError(t, err)
	})
	t.Run("KVSetWithOptions fail", func(t *testing.T) {
		prev := &reacji.List{}
		new := &reacji.List{}
		opt := model.PluginKVSetOptions{
			Atomic:   true,
			OldValue: prev.EncodeToByte(),
		}
		appErr := &model.AppError{}

		api := &plugintest.API{}
		api.On("KVSetWithOptions", keyList, new.EncodeToByte(), opt).Return(false, appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		err := store.Reacji().Update(prev, new)
		assert.Error(t, err)
	})
	t.Run("KVSetWithOptions error", func(t *testing.T) {
		prev := &reacji.List{}
		new := &reacji.List{}
		opt := model.PluginKVSetOptions{
			Atomic:   true,
			OldValue: prev.EncodeToByte(),
		}

		api := &plugintest.API{}
		api.On("KVSetWithOptions", keyList, new.EncodeToByte(), opt).Return(false, nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		err := store.Reacji().Update(prev, new)
		assert.Error(t, err)
	})
}

func TestReacjiStoreForceUpdate(t *testing.T) {
	t.Run("all fine", func(t *testing.T) {
		in := &reacji.List{}

		api := &plugintest.API{}
		api.On("KVSet", keyList, in.EncodeToByte()).Return(nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		err := store.Reacji().ForceUpdate(in)
		assert.NoError(t, err)
	})
	t.Run("KVSet fail", func(t *testing.T) {
		in := &reacji.List{}
		appErr := &model.AppError{}

		api := &plugintest.API{}
		api.On("KVSet", keyList, in.EncodeToByte()).Return(appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		err := store.Reacji().ForceUpdate(in)
		assert.Error(t, err)
	})
}
