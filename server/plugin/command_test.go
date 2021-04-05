package plugin

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kaakaa/mattermost-plugin-reacji/server/store/mockstore"
	"github.com/kaakaa/mattermost-plugin-reacji/server/utils/testutils"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPluginExecuteCommend(t *testing.T) {
	customEmojiName := "custom_emoji"

	for name, test := range map[string]struct { // nolint: govet
		SetupAPI         func(*plugintest.API) *plugintest.API
		SetupHelpers     func(*plugintest.Helpers) *plugintest.Helpers
		SetupStore       func(*mockstore.Store) *mockstore.Store
		Args             *model.CommandArgs
		ShouldError      bool
		ExpectedResponse *model.CommandResponse
	}{
		"error, invalid command": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogError", testutils.GetMockArgumentsWithType("string", 3)...)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
				Command:   "/invalid",
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "invalid command",
			},
		},
		"fine, no args": {
			SetupAPI:     func(api *plugintest.API) *plugintest.API { return api },
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
				Command:   "/reacji",
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: commandHelpMessage,
			},
		},
		"fine, no args with a space": {
			SetupAPI:     func(api *plugintest.API) *plugintest.API { return api },
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
				Command:   "/reacji ",
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: commandHelpMessage,
			},
		},
		"fine, help": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
				Command:   "/reacji help",
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: commandHelpMessage,
			},
		},
		"fine, add with system emoji": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Update", mock.AnythingOfType("*reacji.List"), mock.AnythingOfType("*reacji.List")).Return(nil)
				return s
			},
			Args: &model.CommandArgs{
				UserId:          testutils.GetUserID(),
				ChannelId:       testutils.GetChannelID(),
				Command:         "/reacji add :+1: ~channel",
				ChannelMentions: model.ChannelMentionMap{"channel": testutils.GetChannelID2()},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "add reacjis successfully",
			},
		},
		"fine, add with multiple emojis": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				api.On("GetEmojiByName", customEmojiName).Return(nil, nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Update", mock.AnythingOfType("*reacji.List"), mock.AnythingOfType("*reacji.List")).Return(nil)
				return s
			},
			Args: &model.CommandArgs{
				UserId:          testutils.GetUserID(),
				ChannelId:       testutils.GetChannelID(),
				Command:         fmt.Sprintf("/reacji add :%s: :%s: ~channel", "+1", customEmojiName),
				ChannelMentions: model.ChannelMentionMap{"channel": testutils.GetChannelID2()},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "add reacjis successfully",
			},
		},
		"fine, add with multiple channels": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				api.On("GetEmojiByName", customEmojiName).Return(nil, nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Update", mock.AnythingOfType("*reacji.List"), mock.AnythingOfType("*reacji.List")).Return(nil)
				return s
			},
			Args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
				Command:   fmt.Sprintf("/reacji add :%s: ~channel1 ~channel2", customEmojiName),
				ChannelMentions: model.ChannelMentionMap{
					"channel1": testutils.GetChannelID2(),
					"channel2": testutils.GetChannelID3(),
				},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "add reacjis successfully",
			},
		},
		"fine, add with multiple emojis, multiple channels": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				api.On("GetEmojiByName", customEmojiName).Return(nil, nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Update", mock.AnythingOfType("*reacji.List"), mock.AnythingOfType("*reacji.List")).Return(nil)
				return s
			},
			Args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
				Command:   fmt.Sprintf("/reacji add :%s: :%s: ~channel1 ~channel2", "+1", customEmojiName),
				ChannelMentions: model.ChannelMentionMap{
					"channel1": testutils.GetChannelID2(),
					"channel2": testutils.GetChannelID3(),
				},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "add reacjis successfully",
			},
		},
		"error, add, no emojis": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
				Command:   "/reacji add ~channel1",
				ChannelMentions: model.ChannelMentionMap{
					"channel1": testutils.GetChannelID2(),
				},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "Must specify at least one emoji and at least one channel",
			},
		},
		"error, add, no channels": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args: &model.CommandArgs{
				UserId:          testutils.GetUserID(),
				ChannelId:       testutils.GetChannelID(),
				Command:         "/reacji add :+1:",
				ChannelMentions: model.ChannelMentionMap{},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "Must specify at least one emoji and at least one channel",
			},
		},
		"error, add, updating store fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Update", mock.AnythingOfType("*reacji.List"), mock.AnythingOfType("*reacji.List")).Return(errors.New(""))
				return s
			},
			Args: &model.CommandArgs{
				UserId:    testutils.GetUserID(),
				ChannelId: testutils.GetChannelID(),
				Command:   "/reacji add :+1: ~channel1",
				ChannelMentions: model.ChannelMentionMap{
					"channel1": testutils.GetChannelID2(),
				},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "failed to store reacjis. error=",
			},
		},
		"fine, remove": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Update", mock.AnythingOfType("*reacji.List"), mock.AnythingOfType("*reacji.List")).Return(nil)
				return s
			},
			Args: &model.CommandArgs{
				UserId:          testutils.GetUserID(),
				ChannelId:       testutils.GetChannelID(),
				Command:         fmt.Sprintf("/reacji remove %s", testutils.GetDeleteKey()),
				ChannelMentions: model.ChannelMentionMap{"channel": testutils.GetChannelID2()},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "Reacjis are removed",
			},
		},
		"error, remove, no delete key": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args: &model.CommandArgs{
				UserId:          testutils.GetUserID(),
				ChannelId:       testutils.GetChannelID(),
				Command:         "/reacji remove",
				ChannelMentions: model.ChannelMentionMap{"channel": testutils.GetChannelID2()},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "No delete key",
			},
		},
		"fine, remove-all": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				api.On("GetUser", testutils.GetUserID()).Return(&model.User{Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Update", mock.AnythingOfType("*reacji.List"), mock.AnythingOfType("*reacji.List")).Return(nil)
				return s
			},
			Args: &model.CommandArgs{
				UserId:          testutils.GetUserID(),
				ChannelId:       testutils.GetChannelID(),
				Command:         "/reacji remove-all",
				ChannelMentions: model.ChannelMentionMap{"channel": testutils.GetChannelID2()},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "All reacjis are removed.",
			},
		},
		"fine, remove-all --force": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...)
				api.On("GetUser", testutils.GetUserID()).Return(&model.User{Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("ForceUpdate", mock.AnythingOfType("*reacji.List")).Return(nil)
				return s
			},
			Args: &model.CommandArgs{
				UserId:          testutils.GetUserID(),
				ChannelId:       testutils.GetChannelID(),
				Command:         "/reacji remove-all --force",
				ChannelMentions: model.ChannelMentionMap{"channel": testutils.GetChannelID2()},
			},
			ShouldError: false,
			ExpectedResponse: &model.CommandResponse{
				Text: "All reacjis are removed.",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			a := test.SetupAPI(&plugintest.API{})
			defer a.AssertExpectations(t)
			h := test.SetupHelpers(&plugintest.Helpers{})
			defer h.AssertExpectations(t)
			s := test.SetupStore(&mockstore.Store{})
			defer s.AssertExpectations(t)

			p := setupTestPlugin(a, h, s)
			resp, appErr := p.ExecuteCommand(&plugin.Context{}, test.Args)

			if test.ShouldError {
				assert.NotNil(t, appErr)
			} else {
				assert.Nil(t, appErr)
				assert.Equal(t, test.ExpectedResponse, resp)
			}
		})
	}
}

func TestPluginFindEmojis(t *testing.T) {
	customEmojiName := "custom_emoji"
	customEmojiName2 := "custom_emoji_2"

	var init []string

	for name, test := range map[string]struct {
		SetupAPI     func(*plugintest.API) *plugintest.API
		SetupHelpers func(*plugintest.Helpers) *plugintest.Helpers
		SetupStore   func(*mockstore.Store) *mockstore.Store
		Args         []string
		Expected     []string
	}{
		"find, system emoji": {
			SetupAPI:     func(api *plugintest.API) *plugintest.API { return api },
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args:         []string{":+1:"},
			Expected:     []string{"+1"},
		},
		"fine, custom emoji": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetEmojiByName", customEmojiName).Return(nil, nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args:         []string{fmt.Sprintf(":%s:", customEmojiName)},
			Expected:     []string{customEmojiName},
		},
		"error, GetEmojiByName fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetEmojiByName", customEmojiName).Return(nil, &model.AppError{})
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args:         []string{fmt.Sprintf(":%s:", customEmojiName)},
			Expected:     init,
		},
		"error, GetEmojiByName fails in second times": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetEmojiByName", customEmojiName).Return(nil, nil)
				api.On("GetEmojiByName", customEmojiName2).Return(nil, &model.AppError{})
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			Args:         []string{fmt.Sprintf(":%s:", customEmojiName), fmt.Sprintf(":%s:", customEmojiName2)},
			Expected:     []string{customEmojiName},
		},
	} {
		t.Run(name, func(t *testing.T) {
			a := test.SetupAPI(&plugintest.API{})
			defer a.AssertExpectations(t)
			h := test.SetupHelpers(&plugintest.Helpers{})
			defer h.AssertExpectations(t)
			s := test.SetupStore(&mockstore.Store{})
			defer s.AssertExpectations(t)

			p := setupTestPlugin(a, h, s)

			out := p.findEmojis(test.Args)

			assert.Equal(t, test.Expected, out)
		})
	}
}

func TestPluginRefreshCaches(t *testing.T) {
	for name, test := range map[string]struct { // nolint: govet
		SetupAPI         func(*plugintest.API) *plugintest.API
		SetupHelpers     func(*plugintest.Helpers) *plugintest.Helpers
		SetupStore       func(*mockstore.Store) *mockstore.Store
		UserID           string
		ShouldError      bool
		ExpectedResponse *model.CommandResponse
	}{
		"fine": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetUser", testutils.GetUserID()).Return(&model.User{Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.SharedStore.On("DeleteAll").Return(1, nil)
				return s
			},
			UserID:      testutils.GetUserID(),
			ShouldError: false,
		},
		"error, HasAdminPermission fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetUser", testutils.GetUserID()).Return(nil, &model.AppError{})
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore:   func(s *mockstore.Store) *mockstore.Store { return s },
			UserID:       testutils.GetUserID(),
			ShouldError:  true,
		},
		"error, updating store fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("GetUser", testutils.GetUserID()).Return(&model.User{Roles: model.SYSTEM_ADMIN_ROLE_ID}, nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers { return helpers },
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.SharedStore.On("DeleteAll").Return(0, errors.New(""))
				return s
			},
			UserID:      testutils.GetUserID(),
			ShouldError: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			a := test.SetupAPI(&plugintest.API{})
			defer a.AssertExpectations(t)
			h := test.SetupHelpers(&plugintest.Helpers{})
			defer h.AssertExpectations(t)
			s := test.SetupStore(&mockstore.Store{})
			defer s.AssertExpectations(t)

			p := setupTestPlugin(a, h, s)

			out, appErr := p.refreshCaches(test.UserID)

			if test.ShouldError {
				assert.NotNil(t, appErr)
			} else {
				assert.Nil(t, appErr)
				assert.NotNil(t, out)
			}
		})
	}
}
