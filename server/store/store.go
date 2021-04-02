package store

import "github.com/kaakaa/mattermost-plugin-reacji/server/reacji"

type Store interface {
	Reacji() ReacjiStore
	Shared() SharedStore
}

type ReacjiStore interface {
	Get() (*reacji.List, error)
	Update(prev, new *reacji.List) error
	ForceUpdate(new *reacji.List) error
}

type SharedStore interface {
	Get(postID, toChannelID, deleteKey string) (*reacji.SharedPost, error)
	Set(new *reacji.SharedPost, days int) error
	DeleteAll() (int, error)
}
