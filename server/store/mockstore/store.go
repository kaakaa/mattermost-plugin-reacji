package mockstore

import (
	"github.com/kaakaa/mattermost-plugin-reacji/server/store"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store/mockstore/mocks"
	"github.com/stretchr/testify/mock"
)

type Store struct {
	ReacjiStore mocks.ReacjiStore
	SharedStore mocks.SharedStore
}

func (s *Store) Reacji() store.ReacjiStore { return &s.ReacjiStore }
func (s *Store) Shared() store.SharedStore { return &s.SharedStore }

func (s *Store) AssertExpectations(t mock.TestingT) {
	s.ReacjiStore.AssertExpectations(t)
	s.SharedStore.AssertExpectations(t)
}
