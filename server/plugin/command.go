package plugin

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func (p *Plugin) registerCommand() error {
	return p.API.RegisterCommand(&model.Command{
		Trigger:          "reacji",
		DisplayName:      "Reacji Channeler",
		Description:      "Move post to other channel by attaching reactions",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: add, list, list-all, remove, remove-all, help",
		AutoCompleteHint: "[command]",
		AutocompleteData: createAutoCompleteData(),
	})
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	userID := args.UserId
	FromChannelID := args.ChannelId
	cmdElements := strings.Split(args.Command, " ")

	if len(cmdElements) == 1 || cmdElements[0] != "/"+CommandNameReacji {
		p.API.LogError("invalid command", "command", cmdElements[0])
		return &model.CommandResponse{Text: "invalid command"}, nil
	}
	if len(cmdElements) == 2 && len(cmdElements[1]) == 0 {
		return p.help()
	}
	p.API.LogDebug("execute reacji command", "subcommand", cmdElements[1])
	switch cmdElements[1] {
	case "add":
		emojiNames := p.findEmojis(cmdElements[2:])
		var toChannelIds []string
		for _, id := range args.ChannelMentions {
			toChannelIds = append(toChannelIds, id)
		}
		if len(emojiNames) == 0 || len(toChannelIds) == 0 {
			return &model.CommandResponse{Text: "Must specify at least one emoji and at least one channel"}, nil
		}
		if err := p.storeReacjis(userID, args.TeamId, FromChannelID, emojiNames, toChannelIds); err != nil {
			return &model.CommandResponse{
				Text: fmt.Sprintf("failed to store reacjis. error=%s", err.Error()),
			}, nil
		}
		return &model.CommandResponse{Text: "add reacjis successfully"}, nil
	case "remove":
		if len(cmdElements) == 2 {
			return &model.CommandResponse{Text: "No delete key"}, nil
		}
		return p.remove(userID, cmdElements[2:])
	case "remove-all":
		if len(cmdElements) == 3 && cmdElements[2] == "--force" {
			return p.forceRemoveAll(userID)
		}
		return p.removeAll(userID)
	case "list":
		if len(cmdElements) == 3 && cmdElements[2] == "--all" {
			return p.listAll(userID)
		}
		return p.list(userID, FromChannelID)
	case "help":
		return p.help()
	default:
		return &model.CommandResponse{Text: fmt.Sprintf("invalid subcommand: %s", cmdElements[1])}, nil
	}
}

func (p *Plugin) findEmojis(args []string) []string {
	var ret []string
	re := regexp.MustCompile(`^:[^:]+:$`)
	for _, e := range args {
		text := strings.TrimSpace(e)
		if re.MatchString(text) {
			emojiName := strings.Trim(text, ":")
			if p.isAvailableEmoji(emojiName) {
				ret = append(ret, emojiName)
			}
		}
	}
	return ret
}

func (p *Plugin) isAvailableEmoji(name string) bool {
	// System emoji
	if _, ok := model.SystemEmojis[name]; ok {
		return true
	}
	// Custom emoji
	_, appErr := p.API.GetEmojiByName(name)
	return appErr == nil
}

func (p *Plugin) storeReacjis(userID, teamID, fromChannelID string, emojiNames, toChannelIds []string) error {
	new := p.reacjiList.Clone()
	count := 0
	for _, emoji := range emojiNames {
		for _, to := range toChannelIds {
			if !exists(p.reacjiList, emoji, teamID, to) {
				new.Reacjis = append(new.Reacjis, &reacji.Reacji{
					DeleteKey:     model.NewId(),
					Creator:       userID,
					TeamID:        teamID,
					FromChannelID: fromChannelID,
					ToChannelID:   to,
					EmojiName:     emoji,
				})
				count++
			}
		}
	}
	if count == 0 {
		return errors.New("reacji is already registered")
	}
	if err := p.Store.ReacjiStore.Update(p.reacjiList, new); err != nil {
		return err
	}
	p.reacjiList = new
	p.API.LogDebug("reacjis is updated", "reacjis", fmt.Sprintf("%v", new.Reacjis))
	return nil
}

func exists(list *reacji.List, emoji, teamID, to string) bool {
	for _, reacji := range list.Reacjis {
		if reacji.EmojiName == emoji && reacji.TeamID == teamID && reacji.ToChannelID == to {
			return true
		}
	}
	return false
}

func (p *Plugin) remove(userID string, keys []string) (*model.CommandResponse, *model.AppError) {
	new := &reacji.List{}
	var failed []*reacji.Reacji
	for _, r := range p.reacjiList.Reacjis {
		if include(keys, r.DeleteKey) {
			if b, _ := p.HasAdminPermission(r, userID); b {
				continue
			} else {
				failed = append(failed, r)
				new.Reacjis = append(new.Reacjis, r)
			}
		} else {
			new.Reacjis = append(new.Reacjis, r)
		}
	}
	if err := p.Store.ReacjiStore.Update(p.reacjiList, new); err != nil {
		return &model.CommandResponse{Text: "failed to remove reacjis"}, nil
	}
	p.reacjiList = new
	if len(failed) == 0 {
		return &model.CommandResponse{Text: "Reacjis are removed"}, nil
	}

	var emojis []string
	for _, f := range failed {
		emojis = append(emojis, f.DeleteKey)
	}
	return &model.CommandResponse{Text: fmt.Sprintf("Complete to remove reacjis. However, at least one reacjis encountered error: [%s]\nReacji can be removed by creator or SystemAdministrator.", strings.Join(emojis, ", "))}, nil
}

func include(list []string, key string) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

func (p *Plugin) removeAll(userID string) (*model.CommandResponse, *model.AppError) {
	// TODO: confirm button
	if b, appErr := p.HasAdminPermission(nil, userID); !b {
		appendix := ""
		if appErr != nil {
			appendix = fmt.Sprintf("(%s)", appErr.Error())
		}
		return &model.CommandResponse{
			Text: "Failed to remove emojis due to invalid permission " + appendix,
		}, nil
	}
	new := &reacji.List{}
	if err := p.Store.ReacjiStore.Update(p.reacjiList, new); err != nil {
		return &model.CommandResponse{
			Text: err.Error(),
		}, nil
	}
	p.reacjiList = new
	return &model.CommandResponse{
		Text: "All reacjis are removed.",
	}, nil
}
func (p *Plugin) forceRemoveAll(userID string) (*model.CommandResponse, *model.AppError) {
	// TODO: confirm button
	if b, appErr := p.HasAdminPermission(nil, userID); !b {
		appendix := ""
		if appErr != nil {
			appendix = fmt.Sprintf("(%s)", appErr.Error())
		}
		return &model.CommandResponse{
			Text: "Failed to remove emojis due to nvalid permission " + appendix,
		}, nil
	}

	new := &reacji.List{}
	if err := p.Store.ReacjiStore.ForceUpdate(new); err != nil {
		return &model.CommandResponse{
			Text: err.Error(),
		}, nil
	}
	p.reacjiList = new
	return &model.CommandResponse{
		Text: "All reacjis are removed.",
	}, nil
}

func (p *Plugin) listAll(userID string) (*model.CommandResponse, *model.AppError) {
	table := []string{
		"### All reacjis",
		"",
		"| Emoji | Team | From | To | Creator | DeleteKey | ",
	}
	table = append(table, "|:-----:|:-----|:-----|:---|:--------|:----------|")
	channelCaches := map[string]*model.Channel{}
	for _, r := range p.reacjiList.Reacjis {
		from := fmt.Sprintf("Notfound(ID: %s)", r.FromChannelID)
		if ch, ok := channelCaches[r.FromChannelID]; ok {
			from = fmt.Sprintf("~%s", ch.Name)
		} else {
			fromChannel, appErr := p.API.GetChannel(r.FromChannelID)
			if appErr == nil {
				from = fmt.Sprintf("~%s", fromChannel.Name)
				channelCaches[r.FromChannelID] = fromChannel
			}
		}

		to := fmt.Sprintf("Notfound(ID: %s)", r.ToChannelID)
		if ch, ok := channelCaches[r.ToChannelID]; ok {
			to = fmt.Sprintf("~%s", ch.Name)
		} else {
			toChannel, appErr := p.API.GetChannel(r.ToChannelID)
			if appErr == nil {
				to = fmt.Sprintf("~%s", toChannel.Name)
				channelCaches[r.ToChannelID] = toChannel
			}
		}

		if !p.HasPermissionToPrivateChannel(channelCaches[r.FromChannelID], channelCaches[r.ToChannelID], userID) {
			continue
		}

		teamName := "Unknown"
		team, appErr := p.API.GetTeam(r.TeamID)
		if appErr == nil {
			teamName = team.Name
		} else {
			p.API.LogWarn("failed to get team", "team_id", r.TeamID)
		}

		creator, appErr := p.ConvertUserIDToDisplayName(r.Creator)
		if appErr != nil {
			creator = "Unknown"
		}

		table = append(table, fmt.Sprintf("| :%s: | %s | %s | %s | %s | %s |", r.EmojiName, teamName, from, to, creator, r.DeleteKey))
	}
	if len(table) == 2 {
		return &model.CommandResponse{Text: "There is no Reacji. Add Reacji by `/reacji` command."}, nil
	}
	return &model.CommandResponse{
		Text: strings.Join(table, "\n"),
	}, nil
}

func (p *Plugin) list(userID, channelID string) (*model.CommandResponse, *model.AppError) {
	table := []string{
		"### Reacjis in this channel",
		"",
		"| Emoji | Team | From | To | Creator | DeleteKey | ",
	}
	table = append(table, "|:-----:|:-----|:-----|:---|:--------|:----------|")

	channelCaches := map[string]*model.Channel{}

	fromChannel, appErr := p.API.GetChannel(channelID)
	if appErr != nil {
		return &model.CommandResponse{
			Text: fmt.Sprintf("Failed to get channel by ID: %s", channelID),
		}, nil
	}
	channelCaches[channelID] = fromChannel
	from := fmt.Sprintf("~%s", fromChannel.Name)

	for _, r := range p.reacjiList.Reacjis {
		// Skip if reacji from channel is differ from channel where command is executed
		if r.FromChannelID != channelID {
			continue
		}

		to := fmt.Sprintf("Notfound(ID: %s)", r.ToChannelID)
		if ch, ok := channelCaches[r.ToChannelID]; ok {
			to = fmt.Sprintf("~%s", ch.Name)
		} else {
			toChannel, appErr := p.API.GetChannel(r.ToChannelID)
			if appErr == nil {
				to = fmt.Sprintf("~%s", toChannel.Name)
				channelCaches[r.ToChannelID] = toChannel
			}
		}

		if !p.HasPermissionToPrivateChannel(channelCaches[r.FromChannelID], channelCaches[r.ToChannelID], userID) {
			continue
		}

		teamName := "Unknown"
		team, appErr := p.API.GetTeam(r.TeamID)
		if appErr == nil {
			teamName = team.Name
		} else {
			p.API.LogWarn("failed to get team", "team_id", r.TeamID)
		}

		creator, appErr := p.ConvertUserIDToDisplayName(r.Creator)
		if appErr != nil {
			creator = "Unknown"
		}

		table = append(table, fmt.Sprintf("| :%s: | %s | %s | %s | %s | %s |", r.EmojiName, teamName, from, to, creator, r.DeleteKey))
	}
	if len(table) == 2 {
		return &model.CommandResponse{Text: "There is no Reacji in this channel. Add Reacji by `/reacji` command."}, nil
	}
	return &model.CommandResponse{
		Text: strings.Join(table, "\n"),
	}, nil
}

func (p *Plugin) help() (*model.CommandResponse, *model.AppError) {
	return &model.CommandResponse{
		Text: "Manage Reacjis commands\n" +
			"* `/reacji add :EMOJI: ~CHANNEL`:Register new reacji. If you attach EMOJI to the post, the post will share to CHANNEL.\n" +
			"* `/reacji list [-all]`:List reacjis that is registered in channel. With `-all` list all registered reacjis in this server.\n" +
			"* `/reacji remove [Deletekey...]`: Remove reacjis by DeleteKey (creator or system admin only)\n" +
			"* `/reacji remove-all`: Remove all existing reacjis (system admin only)`\n" +
			"* `/reacji help`: Show help",
	}, nil
}

func createAutoCompleteData() *model.AutocompleteData {
	suggestions := model.NewAutocompleteData("reacji", "[command]", "Available commands: add, list, list-all, remove, remove-all, help")
	suggestions.AddCommand(model.NewAutocompleteData("add", ":EMOJI: ~CHANNEL", "Register new reacji. If you attach EMOJI to the post, the post will share to CHANNEL."))
	suggestions.AddCommand(model.NewAutocompleteData("list", "[-all]", "List reacjis in this channel. With `-all` list all registered reacjis in this server."))
	suggestions.AddCommand(model.NewAutocompleteData("remove", "[DeleteKey...]", "Remove reacji by DeleteKey (creator or system admin only). You can see `DeleteKey` by `/reacji list`"))
	suggestions.AddCommand(model.NewAutocompleteData("remove-all", "", "Remove all reacjis in this server (system admin only)"))
	suggestions.AddCommand(model.NewAutocompleteData("help", "", "Show help"))
	return suggestions
}
