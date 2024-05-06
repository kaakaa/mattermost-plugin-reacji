export interface PluginRegistry {
    registerAppBarComponent(iconUrl: string, action: PluginComponent['action'] | undefined, tooltipText: React.ReactNode, supportedProductIds?: null | string | Array<null | string>, rhsComponent?: PluginComponent['component'] | undefined, rhsTitle?: string)
    registerPostTypeComponent(typeName: string, component: React.ElementType)
    registerRightHandSidebarComponent(component: React.ElementType, title:string)

    // Add more if needed from https://developers.mattermost.com/extend/plugins/webapp/reference
}
