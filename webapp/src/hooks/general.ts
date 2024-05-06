import {useEffect, useState} from 'react';

import { fetchReacjiListByChannelId } from '@/client';

interface Reacji {}

export const useReacjiList = (channelId: string) => {
    const [reacjiList, setReacjiList] = useState<Reacji[]>([]);
    useEffect(() => {
        const fetchReacjiList = async () => {
            const data = await fetchReacjiListByChannelId(channelId);
            setReacjiList(data);
        }
        fetchReacjiList();
    }, [channelId]);
    return reacjiList;
};