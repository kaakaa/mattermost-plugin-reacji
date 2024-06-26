package kvstore

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
	"github.com/kaakaa/mattermost-plugin-reacji/server/utils/testutils"
)

func TestSharedStoreGet(t *testing.T) {
	t.Run("all fine", func(t *testing.T) {
		postID := testutils.GetPostID()
		toChannelID := testutils.GetChannelID()
		deleteKey := testutils.GetDeleteKey()

		key, err := genKey(postID, toChannelID, deleteKey)
		require.NoError(t, err)

		shared := &reacji.SharedPost{}

		api := &plugintest.API{}
		api.On("KVGet", key).Return(shared.EncodeToByte(), nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		out, err := store.Shared().Get(postID, toChannelID, deleteKey)
		assert.NoError(t, err)
		assert.Equal(t, shared, out)
	})
	t.Run("no data", func(t *testing.T) {
		postID := testutils.GetPostID()
		toChannelID := testutils.GetChannelID()
		deleteKey := testutils.GetDeleteKey()

		key, err := genKey(postID, toChannelID, deleteKey)
		require.NoError(t, err)

		api := &plugintest.API{}
		api.On("KVGet", key).Return(nil, nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		out, err := store.Shared().Get(postID, toChannelID, deleteKey)
		assert.NoError(t, err)
		assert.Nil(t, out)
	})
	t.Run("KVGet fail", func(t *testing.T) {
		postID := testutils.GetPostID()
		toChannelID := testutils.GetChannelID()
		deleteKey := testutils.GetDeleteKey()
		appErr := &model.AppError{}

		key, err := genKey(postID, toChannelID, deleteKey)
		require.NoError(t, err)

		api := &plugintest.API{}
		api.On("KVGet", key).Return(nil, appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		out, err := store.Shared().Get(postID, toChannelID, deleteKey)
		assert.Error(t, err)
		assert.Nil(t, out)
	})
	t.Run("invalid data", func(t *testing.T) {
		postID := testutils.GetPostID()
		toChannelID := testutils.GetChannelID()
		deleteKey := testutils.GetDeleteKey()
		appErr := &model.AppError{}

		key, err := genKey(postID, toChannelID, deleteKey)
		require.NoError(t, err)

		api := &plugintest.API{}
		api.On("KVGet", key).Return([]byte{}, appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		out, err := store.Shared().Get(postID, toChannelID, deleteKey)
		assert.Error(t, err)
		assert.Nil(t, out)
	})
}

func TestSharedStoreSet(t *testing.T) {
	t.Run("all fine", func(t *testing.T) {
		shared := &reacji.SharedPost{
			PostID:      testutils.GetPostID(),
			ToChannelID: testutils.GetChannelID(),
			Reacji: reacji.Reacji{
				DeleteKey: testutils.GetDeleteKey(),
			},
		}
		key, err := genKey(shared.PostID, shared.ToChannelID, shared.Reacji.DeleteKey)
		require.NoError(t, err)
		days := 1
		opt := model.PluginKVSetOptions{
			ExpireInSeconds: int64(60 * 60 * 24 * days),
		}

		api := &plugintest.API{}
		api.On("KVSetWithOptions", key, shared.EncodeToByte(), opt).Return(true, nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		err = store.Shared().Set(shared, days)
		assert.NoError(t, err)
	})
	t.Run("KVSetWithOptions fail", func(t *testing.T) {
		shared := &reacji.SharedPost{
			PostID:      testutils.GetPostID(),
			ToChannelID: testutils.GetChannelID(),
			Reacji: reacji.Reacji{
				DeleteKey: testutils.GetDeleteKey(),
			},
		}
		key, err := genKey(shared.PostID, shared.ToChannelID, shared.Reacji.DeleteKey)
		require.NoError(t, err)
		days := 1
		opt := model.PluginKVSetOptions{
			ExpireInSeconds: int64(60 * 60 * 24 * days),
		}

		appErr := &model.AppError{}

		api := &plugintest.API{}
		api.On("KVSetWithOptions", key, shared.EncodeToByte(), opt).Return(false, appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		err = store.Shared().Set(shared, days)
		assert.Error(t, err)
	})
	t.Run("could not set data", func(t *testing.T) {
		shared := &reacji.SharedPost{
			PostID:      testutils.GetPostID(),
			ToChannelID: testutils.GetChannelID(),
			Reacji: reacji.Reacji{
				DeleteKey: testutils.GetDeleteKey(),
			},
		}
		key, err := genKey(shared.PostID, shared.ToChannelID, shared.Reacji.DeleteKey)
		require.NoError(t, err)
		days := 1
		opt := model.PluginKVSetOptions{
			ExpireInSeconds: int64(60 * 60 * 24 * days),
		}

		api := &plugintest.API{}
		api.On("KVSetWithOptions", key, shared.EncodeToByte(), opt).Return(false, nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		err = store.Shared().Set(shared, days)
		assert.Error(t, err)
	})
}

func TestSharedStoreDeleteAll(t *testing.T) {
	t.Run("no data", func(t *testing.T) {
		keys := []string{}

		api := &plugintest.API{}
		api.On("KVList", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(keys, nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		i, err := store.Shared().DeleteAll()
		assert.NoError(t, err)
		assert.Equal(t, len(keys), i)
	})
	t.Run("delete one", func(t *testing.T) {
		keys := []string{"shared-1"}

		api := &plugintest.API{}
		api.On("KVList", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(keys, nil)
		api.On("KVDelete", keys[0]).Return(nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		i, err := store.Shared().DeleteAll()
		assert.NoError(t, err)
		assert.Equal(t, len(keys), i)
	})
	t.Run("delete two", func(t *testing.T) {
		keys := []string{"shared-1", "shared-2"}

		api := &plugintest.API{}
		api.On("KVList", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(keys, nil)
		api.On("KVDelete", keys[0]).Return(nil)
		api.On("KVDelete", keys[1]).Return(nil)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		i, err := store.Shared().DeleteAll()
		assert.NoError(t, err)
		assert.Equal(t, len(keys), i)
	})
	t.Run("KVListWithOptions fail", func(t *testing.T) {
		keys := []string{}
		appErr := &model.AppError{}

		api := &plugintest.API{}
		api.On("KVList", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(keys, appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		i, err := store.Shared().DeleteAll()
		assert.Error(t, err)
		assert.Equal(t, -1, i)
	})
	t.Run("KVDelete fail", func(t *testing.T) {
		keys := []string{"shared-1"}
		appErr := &model.AppError{}

		api := &plugintest.API{}
		api.On("KVList", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(keys, nil)
		api.On("KVDelete", keys[0]).Return(appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		i, err := store.Shared().DeleteAll()
		assert.Error(t, err)
		assert.Equal(t, -1, i)
	})
	t.Run("KVDelete fail in second time", func(t *testing.T) {
		keys := []string{"shared-1", "shared-2"}
		appErr := &model.AppError{}

		api := &plugintest.API{}
		api.On("KVList", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(keys, nil)
		api.On("KVDelete", keys[0]).Return(nil)
		api.On("KVDelete", keys[1]).Return(appErr)
		defer api.AssertExpectations(t)
		store := setupTestStore(api)

		i, err := store.Shared().DeleteAll()
		assert.Error(t, err)
		assert.Equal(t, -1, i)
	})
}

func TestSharedStoreGenKey(t *testing.T) {
	key, err := genKey(testutils.GetPostID(), testutils.GetChannelID(), testutils.GetDeleteKey())
	assert.NoError(t, err)
	assert.Equal(t, "shared-7448bcb561b4c275e4eaf310714e3400", key)
}
