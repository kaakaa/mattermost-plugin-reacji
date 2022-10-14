package testutils

import "github.com/mattermost/mattermost-server/v6/model"

// GetPostID returns a static Post ID.
func GetPostID() string {
	return "post1234567890abcdefghijklm"
}

// GetUserID returns a static User ID.
func GetUserID() string {
	return "user1234567890abcdefghijklm"
}

// GetBotUserID returns a static Bot User ID.
func GetBotUserID() string {
	return "botuser1234567890abcdefghij"
}

// GetTeamID returns a static Team ID.
func GetTeamID() string {
	return "team1234567890abcdefghijklm"
}

// GetChannelID returns a static Channel ID.
func GetChannelID() string {
	return "channel1234567890abcdefghij"
}

// GetChannelID2 returns a static Channel ID.
func GetChannelID2() string {
	return "channel0987654321abcdefghij"
}

// GetChannelID2 returns a static Channel ID.
func GetChannelID3() string {
	return "channelabcdefghijklmnopqrst"
}

// GetDeleteKey returns a static DeleteKey ID.
func GetDeleteKey() string {
	return "deletekeyl1234567890abcdefg"
}

// GetEmojiName returns a static emoji name.
func GetEmojiName() string {
	return "emoji_name"
}

// GetServerConfig return a static server config
func GetServerConfig() *model.Config {
	siteURL := GetSiteURL()
	return &model.Config{
		ServiceSettings: model.ServiceSettings{
			SiteURL: &siteURL,
		},
	}
}

// GetSiteURL return a static site url
func GetSiteURL() string {
	return "https://example.com"
}
