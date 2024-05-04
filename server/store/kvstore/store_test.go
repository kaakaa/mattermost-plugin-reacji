package kvstore

import (
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"
)

func setupTestStore(api plugin.API) *Store {
	kvService := pluginapi.NewClient(api, nil).KV
	store := Store{
		api:         api,
		ReacjiStore: ReacjiStore{api: api},
		SharedStore: SharedStore{api: api, kvService: kvService},
	}
	return &store
}
