package store

import (
	"errors"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const keyList = "reacjis_list"

type Store struct {
	api         plugin.API
	ReacjiStore ReacjiStore
}

func NewStore(api plugin.API) *Store {
	return &Store{
		api:         api,
		ReacjiStore: ReacjiStore{api: api},
	}
}

type ReacjiStore struct {
	api plugin.API
}

func (s *ReacjiStore) Get() (*reacji.List, error) {
	b, err := s.api.KVGet(keyList)
	if err != nil {
		return nil, err
	}
	// b is nil for non-existent
	if b == nil {
		init := &reacji.List{}
		appErr := s.api.KVSet(keyList, init.EncodeToByte())
		if appErr != nil {
			return nil, errors.New("failed to set up kvstore")
		}
		return init, nil
	}

	list := reacji.DecodeListFromByte(b)
	if list == nil {
		return nil, errors.New("failed to decode ReacjisList")
	}
	return list, nil
}

func (s *ReacjiStore) Update(prev, new *reacji.List) error {
	opt := model.PluginKVSetOptions{
		Atomic:   true,
		OldValue: prev.EncodeToByte(),
	}
	ok, err := s.api.KVSetWithOptions(keyList, new.EncodeToByte(), opt)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("failed to store reacji list")
	}
	return nil
}

func (s *ReacjiStore) ForceUpdate(new *reacji.List) error {
	appErr := s.api.KVSet(keyList, new.EncodeToByte())
	if appErr != nil {
		return errors.New(appErr.Error())
	}
	return nil
}
