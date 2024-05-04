package kvstore

import (
	"crypto/md5" // nolint:gosec
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
)

const SharedKeyHeader = "shared-"

type SharedStore struct {
	api       plugin.API
	kvService pluginapi.KVService
}

func (s *SharedStore) Get(postID, toChannelID, deleteKey string) (*reacji.SharedPost, error) {
	key, err := genKey(postID, toChannelID, deleteKey)
	if err != nil {
		return nil, err
	}

	b, appErr := s.api.KVGet(key)
	if appErr != nil {
		return nil, fmt.Errorf("failed to get shared post. %w", appErr)
	}
	if b == nil {
		return nil, nil
	}
	shared := reacji.DecodeSharedPostFromByte(b)
	if shared == nil {
		return nil, errors.New("failed to decode SharedPost")
	}
	return shared, nil
}

func (s *SharedStore) Set(new *reacji.SharedPost, days int) error {
	if days <= 0 {
		return nil
	}

	key, err := genKey(new.PostID, new.ToChannelID, new.Reacji.DeleteKey)
	if err != nil {
		return err
	}

	opt := model.PluginKVSetOptions{
		ExpireInSeconds: int64(60 * 60 * 24 * days),
	}
	ok, appErr := s.api.KVSetWithOptions(key, new.EncodeToByte(), opt)
	if appErr != nil {
		return fmt.Errorf("failed to set shared post. %w", appErr)
	}
	if !ok {
		return errors.New("failed to store shared post")
	}
	return nil
}

func (s *SharedStore) DeleteAll() (int, error) {
	kvListOption := pluginapi.WithPrefix(SharedKeyHeader)
	var count int
	page, perPage := 0, 100
	// TODO: need to test if this loop works fine
	for {
		keys, err := s.kvService.ListKeys(page, perPage, kvListOption)
		if err != nil {
			return -1, err
		}
		for _, k := range keys {
			if err := s.api.KVDelete(k); err != nil {
				return -1, err
			}
			count++
		}
		if len(keys) < perPage {
			break
		}
		page++
	}
	return count, nil
}

func genKey(postID, toChannelID, deleteKey string) (string, error) {
	h := md5.New() // nolint:gosec
	defer h.Reset()
	_, err := h.Write([]byte(fmt.Sprintf("%s-%s-%s", postID, toChannelID, deleteKey)))
	if err != nil {
		return "", err
	}
	v := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s%s", SharedKeyHeader, v), nil
}
