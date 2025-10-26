import * as vscode from 'vscode';
import { ServerManager } from './serverManager';
export declare class ServerStatusItem extends vscode.TreeItem {
    readonly label: string;
    readonly value: string;
    readonly collapsibleState: vscode.TreeItemCollapsibleState;
    readonly iconName?: string | undefined;
    readonly color?: vscode.ThemeColor | undefined;
    constructor(label: string, value: string, collapsibleState: vscode.TreeItemCollapsibleState, iconName?: string | undefined, color?: vscode.ThemeColor | undefined);
}
export declare class ServerStatusProvider implements vscode.TreeDataProvider<ServerStatusItem> {
    private serverManager;
    private _onDidChangeTreeData;
    readonly onDidChangeTreeData: vscode.Event<ServerStatusItem | undefined | null | void>;
    private serverStats;
    private refreshInterval;
    constructor(serverManager: ServerManager);
    refresh(): void;
    getTreeItem(element: ServerStatusItem): vscode.TreeItem;
    getChildren(element?: ServerStatusItem): Thenable<ServerStatusItem[]>;
    private getRootItems;
    private loadServerStats;
    private formatUptime;
    private startAutoRefresh;
    private stopAutoRefresh;
    dispose(): void;
}
//# sourceMappingURL=serverStatusProvider.d.ts.map