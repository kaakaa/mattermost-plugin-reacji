package kvstore

import (
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/plugin"
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
