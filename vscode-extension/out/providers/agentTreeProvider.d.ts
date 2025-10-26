import * as vscode from 'vscode';
import { ServerManager, Agent } from './serverManager';
export declare class AgentItem extends vscode.TreeItem {
    readonly agent: Agent;
    readonly collapsibleState: vscode.TreeItemCollapsibleState;
    constructor(agent: Agent, collapsibleState: vscode.TreeItemCollapsibleState);
}
export declare class AgentCapabilityItem extends vscode.TreeItem {
    readonly capability: string;
    readonly collapsibleState: vscode.TreeItemCollapsibleState;
    constructor(capability: string, collapsibleState: vscode.TreeItemCollapsibleState);
}
export declare class AgentTreeProvider implements vscode.TreeDataProvider<AgentItem | AgentCapabilityItem> {
    private serverManager;
    private _onDidChangeTreeData;
    readonly onDidChangeTreeData: vscode.Event<AgentItem | AgentCapabilityItem | undefined | null | void>;
    private agents;
    private refreshInterval;
    constructor(serverManager: ServerManager);
    refresh(): void;
    getTreeItem(element: AgentItem | AgentCapabilityItem): vscode.TreeItem;
    getChildren(element?: AgentItem | AgentCapabilityItem): Thenable<(AgentItem | AgentCapabilityItem)[]>;
    private getRootItems;
    private loadAgents;
    private startAutoRefresh;
    private stopAutoRefresh;
    getAgentCount(): number;
    getConnectedAgentCount(): number;
    getAgentById(id: string): Agent | undefined;
    dispose(): void;
}
//# sourceMappingURL=agentTreeProvider.d.ts.map