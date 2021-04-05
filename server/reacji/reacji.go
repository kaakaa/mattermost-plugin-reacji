package reacji

import (
	"encoding/json"
)

type List struct {
	Reacjis []*Reacji
}

type Reacji struct {
	DeleteKey     string `json:"delete_key"`
	Creator       string `json:"user_id"`
	TeamID        string `json:"team_id"`
	FromChannelID string `json:"from_channel_id"`
	ToChannelID   string `json:"to_channel_id"`
	EmojiName     string `json:"emoji_name"`
}

type SharedPost struct {
	PostID       string `json:"post_id"`
	ToChannelID  string `json:"to_channel_id"`
	SharedPostID string `json:"shared_post_id"`
	Reacji       Reacji `json:"reacji"`
}

func (l *List) Clone() *List {
	var dst []*Reacji
	for _, r := range l.Reacjis {
		dst = append(dst, r.Clone())
	}
	return &List{Reacjis: dst}
}

func (r *Reacji) Clone() *Reacji {
	return &Reacji{
		DeleteKey:     r.DeleteKey,
		Creator:       r.Creator,
		TeamID:        r.TeamID,
		FromChannelID: r.FromChannelID,
		ToChannelID:   r.ToChannelID,
		EmojiName:     r.EmojiName,
	}
}

func (l *List) EncodeToByte() []byte {
	b, _ := json.Marshal(l)
	return b
}

func DecodeListFromByte(b []byte) *List {
	l := List{}
	if err := json.Unmarshal(b, &l); err != nil {
		return nil
	}
	return &l
}

func (p *SharedPost) EncodeToByte() []byte {
	b, _ := json.Marshal(p)
	return b
}

func DecodeSharedPostFromByte(b []byte) *SharedPost {
	p := SharedPost{}
	if err := json.Unmarshal(b, &p); err != nil {
		return nil
	}
	return &p
}
