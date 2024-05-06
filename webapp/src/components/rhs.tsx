import React from 'react';
import { useSelector } from 'react-redux';
import styled from 'styled-components';

import {GlobalState} from '@mattermost/types/lib/store';
import {Channel} from 'mattermost-redux/types/channels';

import { getCurrentChannelId } from 'mattermost-webapp/webapp/channels/src/packages/mattermost-redux/src/selectors/entities/common';
import { getChannel } from 'mattermost-webapp/webapp/channels/src/packages/mattermost-redux/src/selectors/entities/channels';

import {useReacjiList} from '@/hooks/general';
import RhsRow from '@/components/rhs_row';

const ReactBootstrap = window.ReactBootstrap;

const RhsView = (props: any) => {
    const currentChannelId = useSelector(getCurrentChannelId);
    const channel = useSelector<GlobalState>((state) => getChannel(state, currentChannelId)) as Channel;
    console.log('channel', channel);

    const reacjiList = useReacjiList(currentChannelId);
    const reacjis = reacjiList.map((reacji) => <RhsRow emojiName={reacji.emoji_name} channelId={reacji.to_channel_id}/>);
    return (
        <RhsContainer>
            <RhsTitle>{`Reacjis in ~${channel.display_name}`}</RhsTitle>
            <ReactBootstrap.Table striped bordered hover>
                <thead>
                    <tr>
                        <th>Emoji</th>
                        <th>To Channel</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    {reacjis}
                </tbody>
            </ReactBootstrap.Table>
        </RhsContainer>
    )
}

const RhsTitle = styled.h2``

const RhsContainer = styled.div`
    padding: 5px 20px;
`

export default RhsView;