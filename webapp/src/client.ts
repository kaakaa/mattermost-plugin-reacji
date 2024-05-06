let siteURL = '';
let basePath = '';
let pluginEndpoint = '';

export const setSiteURL = (url: string, pluginId: string) => {
    if (url) {
        basePath = new URL(url).pathname.replace(/\/+$/, '');
        siteURL = url;
    } else {
        basePath = '';
        siteURL = '';
    }

    pluginEndpoint = `${basePath}/plugins/${pluginId}`;
};

export const getIconURL = ():string => {
    return `${pluginEndpoint}/public/logo.dio.png`;
}

export const fetchReacjiListByChannelId = async (channelId: string | null) => {
    let url = `${pluginEndpoint}/api/v1/reacjis`;
    if (channelId) {
        url = `${url}?channel_id=${channelId}`;
    }
    const data = await doGet(url);
    return data;
};

export const doGet = async (url: string) => {
    const response = await fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    });

    console.log('response', response, url);
    if (!response.ok) {
        throw new Error(`Failed to fetch: ${response.statusText}`);
    }

    return response.json();
};
