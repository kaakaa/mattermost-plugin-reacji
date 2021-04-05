package reacji

import (
	"testing"

	"github.com/kaakaa/mattermost-plugin-reacji/server/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestReacjiListClone(t *testing.T) {
	r := []*Reacji{getTestReacji()}
	in := &List{
		Reacjis: r,
	}
	out := in.Clone()
	assert.Equal(t, in, out)
}

func TestDecodeListFromByte(t *testing.T) {
	t.Run("fine", func(t *testing.T) {
		in := &List{Reacjis: []*Reacji{getTestReacji()}}
		out := DecodeListFromByte(in.EncodeToByte())
		assert.Equal(t, in, out)
	})
	t.Run("fail", func(t *testing.T) {
		out := DecodeListFromByte([]byte{})
		assert.Nil(t, out)
	})
}

func TestSharedPostEncodeToByte(t *testing.T) {
	in := &SharedPost{
		PostID:       testutils.GetPostID(),
		ToChannelID:  testutils.GetChannelID(),
		SharedPostID: testutils.GetPostID(),
		Reacji:       *getTestReacji(),
	}
	out := in.EncodeToByte()
	assert.Equal(t, in.EncodeToByte(), out)
}

func TestDecodeSharedPostFromByte(t *testing.T) {
	t.Run("fine", func(t *testing.T) {
		in := &SharedPost{
			PostID:       testutils.GetPostID(),
			ToChannelID:  testutils.GetChannelID(),
			SharedPostID: testutils.GetPostID(),
			Reacji:       *getTestReacji(),
		}
		out := DecodeSharedPostFromByte(in.EncodeToByte())
		assert.Equal(t, in, out)
	})
	t.Run("fail", func(t *testing.T) {
		out := DecodeSharedPostFromByte([]byte{})
		assert.Nil(t, out)
	})
}

// getTestReacji returns a static reacji
func getTestReacji() *Reacji {
	return &Reacji{
		DeleteKey:     testutils.GetDeleteKey(),
		Creator:       testutils.GetUserID(),
		TeamID:        testutils.GetTeamID(),
		FromChannelID: testutils.GetChannelID(),
		ToChannelID:   testutils.GetChannelID2(),
		EmojiName:     testutils.GetEmojiName(),
	}
}
