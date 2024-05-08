import {Store, Action} from 'redux';

import {GlobalState} from '@mattermost/types/store';

import manifest from '@/manifest';

import {PluginRegistry} from '@/types/mattermost-webapp';

import {getConfig} from 'mattermost-redux/selectors/entities/general';

import {getIconURL, setSiteURL} from '@/client';
import RhsView from '@/components/rhs';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        setSiteURL(getConfig(store.getState())?.SiteURL || '', manifest.id);

        const {toggleRHSPlugin} = registry.registerRightHandSidebarComponent(RhsView, 'Reacji List');

        if (registry.registerAppBarComponent) {
            registry.registerAppBarComponent(getIconURL(), () => store.dispatch(toggleRHSPlugin), 'Show Reacjis in current channel');
        }
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
