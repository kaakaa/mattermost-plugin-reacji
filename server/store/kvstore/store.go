package kvstore

import (
	"github.com/kaakaa/mattermost-plugin-reacji/server/store"
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/plugin"
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
