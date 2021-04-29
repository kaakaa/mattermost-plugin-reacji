package plugin

import (
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store/kvstore"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store/mockstore"
	"github.com/kaakaa/mattermost-plugin-reacji/server/utils/testutils"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestPlugin(api *plugintest.API, helpers *plugintest.Helpers, store *mockstore.Store) *Plugin {
	p := &Plugin{
		ServerConfig: testutils.GetServerConfig(),
	}
	p.setConfiguration(&configuration{
		AllowDuplicateSharing:  true,
		DaysToKeepSharedRecord: 1,
		MaxReacjis:             30,
	})

	p.SetAPI(api)
	p.SetHelpers(helpers)
	p.botUserID = testutils.GetBotUserID()
	p.reacjiList = &reacji.List{Reacjis: []*reacji.Reacji{getTestReacji()}}
	p.Store = store

	return p
}

// getTestReacji returns a static reacji
func getTestReacji() *reacji.Reacji {
	return &reacji.Reacji{
		DeleteKey:     testutils.GetDeleteKey(),
		Creator:       testutils.GetUserID(),
		TeamID:        testutils.GetTeamID(),
		FromChannelID: testutils.GetChannelID(),
		ToChannelID:   testutils.GetChannelID2(),
		EmojiName:     testutils.GetEmojiName(),
	}
}

func TestPluginOnActivate(t *testing.T) {
	for name, test := range map[string]struct {
		SetupAPI     func(*plugintest.API) *plugintest.API
		SetupHelpers func(*plugintest.Helpers) *plugintest.Helpers
		SetupStore   func(*mockstore.Store) *mockstore.Store
		ShouldError  bool
	}{
		"fine": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 1)...).Return(nil)
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...).Return(nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(testutils.GetBotUserID(), nil)
				return helpers
			},
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Get").Return(&reacji.List{}, nil)
				return s
			},
			ShouldError: false,
		},
		"error, Helpers.EnsureBot fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 1)...).Return(nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return("", &model.AppError{})
				return helpers
			},
			SetupStore:  func(s *mockstore.Store) *mockstore.Store { return s },
			ShouldError: true,
		},
		"error, getting reacji from store fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 1)...).Return(nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(testutils.GetBotUserID(), nil)
				return helpers
			},
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Get").Return(nil, errors.New(""))
				return s
			},
			ShouldError: true,
		},
		"error, RegisterCommand fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(&model.AppError{})
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 1)...).Return(nil)
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...).Return(nil)
				return api
			},
			SetupHelpers: func(helpers *plugintest.Helpers) *plugintest.Helpers {
				helpers.On("EnsureBot", mock.AnythingOfType("*model.Bot"), mock.AnythingOfType("plugin.EnsureBotOption")).Return(testutils.GetBotUserID(), nil)
				return helpers
			},
			SetupStore: func(s *mockstore.Store) *mockstore.Store {
				s.ReacjiStore.On("Get").Return(&reacji.List{}, nil)
				return s
			},
			ShouldError: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			a := test.SetupAPI(&plugintest.API{})
			defer a.AssertExpectations(t)
			h := test.SetupHelpers(&plugintest.Helpers{})
			defer h.AssertExpectations(t)
			s := test.SetupStore(&mockstore.Store{})
			defer s.AssertExpectations(t)

			patch := monkey.Patch(kvstore.NewStore, func(plugin.API, plugin.Helpers) store.Store {
				return s
			})
			defer patch.Unpatch()

			p := setupTestPlugin(a, h, s)
			err := p.OnActivate()

			if test.ShouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPluginOnDeactivate(t *testing.T) {
	t.Run("fine", func(t *testing.T) {
		a := &plugintest.API{}
		a.On("UnregisterCommand", "", CommandNameReacji).Return(nil)
		defer a.AssertExpectations(t)
		h := &plugintest.Helpers{}
		defer h.AssertExpectations(t)
		s := &mockstore.Store{}
		defer s.AssertExpectations(t)

		p := setupTestPlugin(a, h, s)
		err := p.OnDeactivate()

		assert.NoError(t, err)
	})
	t.Run("error, UnregisterCommand fails", func(t *testing.T) {
		a := &plugintest.API{}
		a.On("UnregisterCommand", "", CommandNameReacji).Return(errors.New(""))
		defer a.AssertExpectations(t)
		h := &plugintest.Helpers{}
		defer h.AssertExpectations(t)
		s := &mockstore.Store{}
		defer s.AssertExpectations(t)

		p := setupTestPlugin(a, h, s)
		err := p.OnDeactivate()

		assert.Error(t, err)
	})
}

func TestPluginHasPermissionToChannel(t *testing.T) {
	userID := testutils.GetUserID()
	channelID := testutils.GetChannelID()

	for name, test := range map[string]struct {
		SetupAPI func(*plugintest.API) *plugintest.API
		Channel  *model.Channel
		UserID   string
		Expected bool
	}{
		"fine, public channel": {
			SetupAPI: func(api *plugintest.API) *plugintest.API { return api },
			Channel:  &model.Channel{Id: channelID, Type: model.CHANNEL_OPEN},
			UserID:   userID,
			Expected: true,
		},
		"fine, private channel with permission": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("HasPermissionToChannel", userID, channelID, model.PERMISSION_READ_CHANNEL).Return(true)
				return api
			},
			Channel:  &model.Channel{Id: channelID, Type: model.CHANNEL_PRIVATE},
			UserID:   userID,
			Expected: true,
		},
		"fine, private channel without permission": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("HasPermissionToChannel", userID, channelID, model.PERMISSION_READ_CHANNEL).Return(false)
				return api
			},
			Channel:  &model.Channel{Id: channelID, Type: model.CHANNEL_PRIVATE},
			UserID:   userID,
			Expected: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			a := test.SetupAPI(&plugintest.API{})
			defer a.AssertExpectations(t)
			h := &plugintest.Helpers{}
			defer h.AssertExpectations(t)
			s := &mockstore.Store{}
			defer s.AssertExpectations(t)

			p := setupTestPlugin(a, h, s)
			actual := p.HasPermissionToChannel(test.Channel, test.UserID)

			assert.Equal(t, test.Expected, actual)
		})
	}
}
