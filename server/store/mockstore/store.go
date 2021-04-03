package mockstore

import (
	"github.com/kaakaa/mattermost-plugin-reacji/server/store"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store/mockstore/mocks"
	"github.com/stretchr/testify/mock"
)

type Store struct {
	ReacjiStore mocks.ReacjiStore
}

func (s *Store) Reacji() store.ReacjiStore {
	return &s.ReacjiStore
}

func (s *Store) AssertExpectations(t mock.TestingT) {
	s.ReacjiStore.AssertExpectations(t)
}
