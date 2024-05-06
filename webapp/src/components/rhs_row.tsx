import React from 'react';
import { useSelector } from 'react-redux';
import styled from 'styled-components';

import {GlobalState} from '@mattermost/types/lib/store';
import {Channel} from 'mattermost-redux/types/channels';

import { getChannel, getChannelsNameMapInCurrentTeam } from 'mattermost-webapp/webapp/channels/src/packages/mattermost-redux/src/selectors/entities/channels';
import { getCurrentTeam } from 'mattermost-webapp/webapp/channels/src/packages/mattermost-redux/src/selectors/entities/teams';

import {getSiteURL} from '@/client';

const PostUtils = window.PostUtils;

type RhsRowProps = {
    emojiName: string;
    channelId: string;
}

const RhsRow = ({emojiName, channelId}: RhsRowProps) => {
    const siteURL = getSiteURL();
    const currentTeam = useSelector(getCurrentTeam);
    const channelNamesMap = useSelector<GlobalState>((state) => getChannelsNameMapInCurrentTeam(state)) as Record<string, Channel>;

    const c = PostUtils.messageHtmlToComponent(
        PostUtils.formatText(`:${emojiName}:`),
        true
    );

    const channel = useSelector<GlobalState>((state) => getChannel(state, channelId)) as Channel;
    const ch = PostUtils.messageHtmlToComponent(
        PostUtils.formatText(`~${channel.name}`, {siteURL: siteURL, channelNamesMap: channelNamesMap, team: currentTeam}),
        true
    );

    return (
        <tr>
            <RhsCell>{c}</RhsCell>
            <RhsCell>{ch}</RhsCell>
            <RhsCell><DeleteButton>{'DELETE'}</DeleteButton></RhsCell>
        </tr>
    )
};

const RhsCell = styled.td`
    vertical-align: middle;
`

const DeleteButton = styled.button`
    color: #AA043D;
    border: 2px solid #AA043D;
    border-radius: 8px;
    transition-duration: 0.4s;
    padding: 5px 10px;
    &:hover {
        background-color: #AA043D;
        color: white;
    }
`

export default RhsRow;