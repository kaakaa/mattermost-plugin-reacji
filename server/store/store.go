package store

import "github.com/kaakaa/mattermost-plugin-reacji/server/reacji"

type Store interface {
	Reacji() ReacjiStore
}

type ReacjiStore interface {
	Get() (*reacji.List, error)
	Update(prev, new *reacji.List) error
	ForceUpdate(new *reacji.List) error
}
