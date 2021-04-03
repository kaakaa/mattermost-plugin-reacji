package kvstore

import "github.com/mattermost/mattermost-server/v5/plugin"

func setupTestStore(api plugin.API, helpers plugin.Helpers) *Store {
	store := Store{
		api:         api,
		ReacjiStore: ReacjiStore{api: api},
		SharedStore: SharedStore{api: api, helpers: helpers},
	}
	return &store
}
