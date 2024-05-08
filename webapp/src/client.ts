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

export const getSiteURL = ():string => {
    return siteURL;
};

export const getIconURL = ():string => {
    return `${pluginEndpoint}/public/logo.dio.png`;
};

export const fetchReacjiListByChannelId = async (channelId: string | null) => {
    let url = `${pluginEndpoint}/api/v1/reacjis`;
    if (channelId) {
        url = `${url}?channel_id=${channelId}`;
    }
    const data = await doGet(url);
    return data;
};

export const openDeleteReacjiConfirmationDialog = async (deleteKey: string) => {
    // @ts-ignore
    window.openInteractiveDialog({ // ref: https://github.com/mattermost/mattermost-webapp/pull/2838
        url: `${pluginEndpoint}/api/v1/reacjis/${deleteKey}/confirm`,
        dialog: {
            callback_id: deleteKey,
            title: 'Delete Reacji',
            introduction_text: 'NOTE: If you delete a reacji from RHS, this deletion will not fire re-rendering RHS component. Please re-open RHS to see the updated list.',
            elements: [],
            submit_label: 'Confirm',
            state: 'delete_confirmation',
        },
    });
};

const doGet = async (url: string, headers?: Record<string, string>) => {
    return doRequest(url, 'GET', {'Content-Type': 'application/json'});
};

const doRequest = async (url: string, method: string, headers?: Record<string, string>) => {
    const response = await fetch(url, {method, headers});

    if (!response.ok) {
        throw new Error(`Failed to fetch: ${response.statusText}`);
    }

    return response.json();
};
