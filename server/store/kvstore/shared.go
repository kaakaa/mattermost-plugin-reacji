package kvstore

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const SharedKeyHeader = "shared"

type SharedStore struct {
	api plugin.API
}

func (s *SharedStore) Get(postID, toChannelID, deleteKey string) (*reacji.SharedPost, error) {
	key := genKey(postID, toChannelID, deleteKey)
	b, err := s.api.KVGet(key)
	if err != nil {
		return nil, err
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
	key := genKey(new.PostID, new.ToChannelID, new.Reacji.DeleteKey)

	opt := model.PluginKVSetOptions{
		ExpireInSeconds: int64(time.Hour * 24 * time.Duration(days)),
	}
	ok, err := s.api.KVSetWithOptions(key, new.EncodeToByte(), opt)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("failed to store shared post")
	}
	return nil
}

func genKey(postID, toChannelID, deleteKey string) string {
	h := md5.New()
	defer h.Reset()
	h.Write([]byte(fmt.Sprintf("%s-%s-%s-%s", SharedKeyHeader, postID, toChannelID, deleteKey)))
	return hex.EncodeToString(h.Sum(nil))
}
