import React from 'react';
import {useSelector} from 'react-redux';
import styled from 'styled-components';

import {Reacji} from '@/types/types';

import {GlobalState} from '@mattermost/types/store';
import {Channel} from '@mattermost/types/channels';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';

import {useReacjiList} from '@/hooks/general';
import RhsRow from '@/components/rhs_row';

const ReactBootstrap = window.ReactBootstrap;

const RhsView = () => {
    const currentChannelId = useSelector(getCurrentChannelId);
    const channel = useSelector<GlobalState>((state) => getChannel(state, currentChannelId)) as Channel;

    const reacjiList = useReacjiList(currentChannelId);
    const reacjis = reacjiList.map((reacji: Reacji) => {
        return (
            <RhsRow
                key={reacji.delete_key}
                emojiName={reacji.emoji_name}
                channelId={reacji.to_channel_id}
            />
        );
    });
    return (
        <RhsContainer>
            <RhsTitle>{`Reacjis in ~${channel.display_name}`}</RhsTitle>
            <ReactBootstrap.Table
                striped={true}
                bordered={true}
                hover={true}
            >
                <thead>
                    <tr>
                        <th>{'Emoji'}</th>
                        <th>{'To Channel'}</th>
                        <th/>
                    </tr>
                </thead>
                <tbody>
                    {reacjis}
                </tbody>
            </ReactBootstrap.Table>
        </RhsContainer>
    );
};

const RhsTitle = styled.h2``;

const RhsContainer = styled.div`
    padding: 5px 20px;
`;

export default RhsView;