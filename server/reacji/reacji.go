package reacji

import (
	"encoding/json"
)

type ReacjiList struct {
	Reacjis []*Reacji
}

type Reacji struct {
	DeleteKey     string `json:"delete_key"`
	Creator       string `json:"user_id"`
	TeamId        string `json:"team_id"`
	FromChannelId string `json:"from_channel_id"`
	ToChannelId   string `json:"to_channel_id"`
	EmojiName     string `json:"emoji_name"`
}

func (l *ReacjiList) Clone() *ReacjiList {
	var dst []*Reacji
	for _, r := range l.Reacjis {
		dst = append(dst, r.Clone())
	}
	return &ReacjiList{Reacjis: dst}
}

func (r *Reacji) Clone() *Reacji {
	return &Reacji{
		DeleteKey:     r.DeleteKey,
		Creator:       r.Creator,
		TeamId:        r.TeamId,
		FromChannelId: r.FromChannelId,
		ToChannelId:   r.ToChannelId,
		EmojiName:     r.EmojiName,
	}
}

func (l *ReacjiList) EncodeToByte() []byte {
	b, _ := json.Marshal(l)
	return b
}

func DecodeReacjiListFromByte(b []byte) *ReacjiList {
	l := ReacjiList{}
	if err := json.Unmarshal(b, &l); err != nil {
		return nil
	}
	return &l
}
