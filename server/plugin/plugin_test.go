package plugin

import (
	"errors"
	"reflect"
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/undefinedlabs/go-mpatch"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store/kvstore"
	"github.com/kaakaa/mattermost-plugin-reacji/server/store/mockstore"
	"github.com/kaakaa/mattermost-plugin-reacji/server/utils/testutils"
)

func setupTestPlugin(api *plugintest.API, store *mockstore.Store) *Plugin {
	p := &Plugin{
		ServerConfig: testutils.GetServerConfig(),
	}
	p.setConfiguration(&configuration{
		AllowDuplicateSharing:  true,
		DaysToKeepSharedRecord: 1,
		MaxReacjis:             30,
	})

	p.SetAPI(api)
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
		SetupAPI       func(*plugintest.API) *plugintest.API
		SetupPluginAPI func(*pluginapi.Client) (*pluginapi.Client, []*mpatch.Patch)
		SetupStore     func(*mockstore.Store) *mockstore.Store
		ShouldError    bool
	}{
		// Disable this test because mocking store.Store is not work fine.
		"fine": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("RegisterCommand", mock.AnythingOfType("*model.Command")).Return(nil)
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 1)...).Return(nil)
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 3)...).Return(nil)
				return api
			},
			SetupPluginAPI: func(client *pluginapi.Client) (*pluginapi.Client, []*mpatch.Patch) {
				p, err := mpatch.PatchInstanceMethodByName(reflect.TypeOf(client.Bot), "EnsureBot", func(*pluginapi.BotService, *model.Bot, ...pluginapi.EnsureBotOption) (string, error) {
					return testutils.GetBotUserID(), nil
				})
				require.NoError(t, err)
				return client, []*mpatch.Patch{p}
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
			SetupPluginAPI: func(client *pluginapi.Client) (*pluginapi.Client, []*mpatch.Patch) {
				p, err := mpatch.PatchInstanceMethodByName(reflect.TypeOf(client.Bot), "EnsureBot", func(*pluginapi.BotService, *model.Bot, ...pluginapi.EnsureBotOption) (string, error) {
					return "", errors.New("")
				})
				require.NoError(t, err)
				return client, []*mpatch.Patch{p}
			},
			SetupStore:  func(s *mockstore.Store) *mockstore.Store { return s },
			ShouldError: true,
		},
		"error, getting reacji from store fails": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("LogDebug", testutils.GetMockArgumentsWithType("string", 1)...).Return(nil)
				return api
			},
			SetupPluginAPI: func(client *pluginapi.Client) (*pluginapi.Client, []*mpatch.Patch) {
				p, err := mpatch.PatchInstanceMethodByName(reflect.TypeOf(client.Bot), "EnsureBot", func(*pluginapi.BotService, *model.Bot, ...pluginapi.EnsureBotOption) (string, error) {
					return testutils.GetBotUserID(), nil
				})
				require.NoError(t, err)
				return client, []*mpatch.Patch{p}
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
			SetupPluginAPI: func(client *pluginapi.Client) (*pluginapi.Client, []*mpatch.Patch) {
				p, err := mpatch.PatchInstanceMethodByName(reflect.TypeOf(client.Bot), "EnsureBot", func(*pluginapi.BotService, *model.Bot, ...pluginapi.EnsureBotOption) (string, error) {
					return testutils.GetBotUserID(), nil
				})
				require.NoError(t, err)
				return client, []*mpatch.Patch{p}
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
			s := test.SetupStore(&mockstore.Store{})
			defer s.AssertExpectations(t)

			patch1, err := mpatch.PatchMethod(
				kvstore.NewStore,
				func(plugin.API, pluginapi.KVService) store.Store { return s },
			)
			require.NoError(t, err)
			defer func() { require.NoError(t, patch1.Unpatch()) }()

			// Setup pluginapi client
			mClient := pluginapi.NewClient(a, &plugintest.Driver{})
			patch2, err := mpatch.PatchMethod(
				pluginapi.NewClient,
				func(plugin.API, plugin.Driver) *pluginapi.Client { return mClient },
			)
			require.NoError(t, err)
			defer func() { require.NoError(t, patch2.Unpatch()) }()

			if test.SetupPluginAPI != nil {
				_, patches := test.SetupPluginAPI(mClient)
				t.Cleanup(func() {
					for _, p := range patches {
						require.NoError(t, p.Unpatch())
					}
				})
			}

			p := setupTestPlugin(a, s)
			err = p.OnActivate()

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
		s := &mockstore.Store{}
		defer s.AssertExpectations(t)

		p := setupTestPlugin(a, s)
		err := p.OnDeactivate()

		assert.NoError(t, err)
	})
	t.Run("error, UnregisterCommand fails", func(t *testing.T) {
		a := &plugintest.API{}
		a.On("UnregisterCommand", "", CommandNameReacji).Return(errors.New(""))
		defer a.AssertExpectations(t)
		s := &mockstore.Store{}
		defer s.AssertExpectations(t)

		p := setupTestPlugin(a, s)
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
			Channel:  &model.Channel{Id: channelID, Type: model.ChannelTypeOpen},
			UserID:   userID,
			Expected: true,
		},
		"fine, private channel with permission": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("HasPermissionToChannel", userID, channelID, model.PermissionReadChannel).Return(true)
				return api
			},
			Channel:  &model.Channel{Id: channelID, Type: model.ChannelTypePrivate},
			UserID:   userID,
			Expected: true,
		},
		"fine, private channel without permission": {
			SetupAPI: func(api *plugintest.API) *plugintest.API {
				api.On("HasPermissionToChannel", userID, channelID, model.PermissionReadChannel).Return(false)
				return api
			},
			Channel:  &model.Channel{Id: channelID, Type: model.ChannelTypePrivate},
			UserID:   userID,
			Expected: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			a := test.SetupAPI(&plugintest.API{})
			defer a.AssertExpectations(t)
			s := &mockstore.Store{}
			defer s.AssertExpectations(t)

			p := setupTestPlugin(a, s)
			actual := p.HasPermissionToChannel(test.Channel, test.UserID)

			assert.Equal(t, test.Expected, actual)
		})
	}
}
