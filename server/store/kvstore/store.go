package kvstore

import (
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/kaakaa/mattermost-plugin-reacji/server/store"
)

type Store struct {
	api         plugin.API
	ReacjiStore ReacjiStore
	SharedStore SharedStore
}

func NewStore(api plugin.API, kvService pluginapi.KVService) store.Store {
	return &Store{
		api:         api,
		ReacjiStore: ReacjiStore{api: api},
		SharedStore: SharedStore{api: api, kvService: kvService},
	}
}

func (s *Store) Reacji() store.ReacjiStore {
	return &s.ReacjiStore
}

func (s *Store) Shared() store.SharedStore {
	return &s.SharedStore
}
