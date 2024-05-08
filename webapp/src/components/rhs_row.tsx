import React from 'react';
import {useSelector} from 'react-redux';
import styled from 'styled-components';

import {GlobalState} from '@mattermost/types/store';
import {Channel} from '@mattermost/types/channels';

import {getChannel, getChannelsNameMapInCurrentTeam} from 'mattermost-redux/selectors/entities/channels';
import {getCurrentTeam} from 'mattermost-redux/selectors/entities/teams';

import {Reacji} from '@/types/types';
import {openDeleteReacjiConfirmationDialog, getSiteURL} from '@/client';

// @ts-ignore
const PostUtils = window.PostUtils;

type RhsRowProps = {
    reacji: Reacji;
};

const RhsRow = ({reacji}: RhsRowProps) => {
    const {emoji_name, to_channel_id, delete_key} = reacji;

    const siteURL = getSiteURL();
    const currentTeam = useSelector(getCurrentTeam);
    const channelNamesMap = useSelector<GlobalState>((state) => getChannelsNameMapInCurrentTeam(state)) as Record<string, Channel>;

    const c = PostUtils.messageHtmlToComponent(
        PostUtils.formatText(`:${emoji_name}:`),
        true,
    );

    const channel = useSelector<GlobalState>((state) => getChannel(state, to_channel_id)) as Channel;
    const ch = PostUtils.messageHtmlToComponent(
        PostUtils.formatText(`~${channel.name}`, {siteURL, channelNamesMap, team: currentTeam}),
        true,
    );

    return (
        <tr>
            <RhsCell>{c}</RhsCell>
            <RhsCell>{ch}</RhsCell>
            <RhsCell
                onClick={() => openDeleteReacjiConfirmationDialog(delete_key)}
            >
                <DeleteButton>{'DELETE'}</DeleteButton>
            </RhsCell>
        </tr>
    );
};

const RhsCell = styled.td`
    vertical-align: middle;
`;

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
`;

export default RhsRow;