import {useEffect, useState} from 'react';

import {Reacji} from '@/types/types';
import {fetchReacjiListByChannelId} from '@/client';

export const useReacjiList = (channelId: string) => {
    const [reacjiList, setReacjiList] = useState<Reacji[]>([]);
    useEffect(() => {
        const fetchReacjiList = async () => {
            const data = await fetchReacjiListByChannelId(channelId);
            setReacjiList(data);
        };
        fetchReacjiList();
    }, [channelId]);
    return reacjiList;
};
