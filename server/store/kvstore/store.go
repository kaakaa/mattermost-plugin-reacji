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

func NewStore(api plugin.API) store.Store {
	return &Store{
		api:         api,
		ReacjiStore: ReacjiStore{api: api},
		SharedStore: SharedStore{api: api},
	}
}

func (s *Store) Reacji() store.ReacjiStore {
	return &s.ReacjiStore
}

func (s *Store) Shared() store.SharedStore {
	return &s.SharedStore
}
