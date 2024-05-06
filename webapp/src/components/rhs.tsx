import React from 'react';
import { useSelector } from 'react-redux';

import { getCurrentChannelId } from 'mattermost-webapp/webapp/channels/src/packages/mattermost-redux/src/selectors/entities/common';

import {useReacjiList} from '@/hooks/general';

const ReactBootstrap = window.ReactBootstrap;

const RhsView = (props: any) => {
    const channelId = useSelector(getCurrentChannelId);
    const reacjiList = useReacjiList(channelId);

    const reacjis = reacjiList.map((reacji) => {
        return (
            <ReactBootstrap.Row>
                <ReactBootstrap.Col>{reacji.emoji_name}</ReactBootstrap.Col>
                <ReactBootstrap.Col>{reacji.to_channel_id}</ReactBootstrap.Col>
                <ReactBootstrap.Col>{'DELETE'}</ReactBootstrap.Col>
            </ReactBootstrap.Row>
        )
    });
    return (
        <>
            <h1>Reacji RHS View</h1>
            <div>
                {reacjis}
            </div>
        </>
    )
}

export default RhsView;