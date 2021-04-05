package kvstore

import (
	"github.com/kaakaa/mattermost-plugin-reacji/server/store"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

type Store struct {
	api         plugin.API
	ReacjiStore ReacjiStore
	SharedStore SharedStore
}

func NewStore(api plugin.API, helpers plugin.Helpers) store.Store {
	return &Store{
		api:         api,
		ReacjiStore: ReacjiStore{api: api},
		SharedStore: SharedStore{api: api, helpers: helpers},
	}
}

func (s *Store) Reacji() store.ReacjiStore {
	return &s.ReacjiStore
}

func (s *Store) Shared() store.SharedStore {
	return &s.SharedStore
}
