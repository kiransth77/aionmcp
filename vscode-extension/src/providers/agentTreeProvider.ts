import * as vscode from 'vscode';
import { ServerManager, Agent } from './serverManager';

export class AgentItem extends vscode.TreeItem {
    constructor(
        public readonly agent: Agent,
        public readonly collapsibleState: vscode.TreeItemCollapsibleState
    ) {
        super(agent.name || agent.id, collapsibleState);
        
        this.tooltip = `Agent: ${agent.id}\nStatus: ${agent.status}\nCapabilities: ${agent.capabilities.join(', ')}`;
        this.description = agent.status;
        this.contextValue = 'agent';
        
        // Set icon and color based on status
        if (agent.status === 'connected') {
            this.iconPath = new vscode.ThemeIcon('account', new vscode.ThemeColor('charts.green'));
        } else {
            this.iconPath = new vscode.ThemeIcon('account', new vscode.ThemeColor('charts.red'));
        }
        
        // Show last seen time for disconnected agents
        if (agent.status === 'disconnected' && agent.lastSeen) {
            const timeDiff = Date.now() - agent.lastSeen.getTime();
            const minutes = Math.floor(timeDiff / 60000);
            const hours = Math.floor(minutes / 60);
            
            if (hours > 0) {
                this.description = `${agent.status} (${hours}h ago)`;
            } else if (minutes > 0) {
                this.description = `${agent.status} (${minutes}m ago)`;
            } else {
                this.description = `${agent.status} (just now)`;
            }
        }
    }
}

export class AgentCapabilityItem extends vscode.TreeItem {
    constructor(
        public readonly capability: string,
        public readonly collapsibleState: vscode.TreeItemCollapsibleState
    ) {
        super(capability, collapsibleState);
        
        this.tooltip = `Capability: ${capability}`;
        this.contextValue = 'agentCapability';
        this.iconPath = new vscode.ThemeIcon('symbol-property');
    }
}

export class AgentTreeProvider implements vscode.TreeDataProvider<AgentItem | AgentCapabilityItem> {
    private _onDidChangeTreeData: vscode.EventEmitter<AgentItem | AgentCapabilityItem | undefined | null | void> = new vscode.EventEmitter<AgentItem | AgentCapabilityItem | undefined | null | void>();
    readonly onDidChangeTreeData: vscode.Event<AgentItem | AgentCapabilityItem | undefined | null | void> = this._onDidChangeTreeData.event;
    
    private agents: Agent[] = [];
    private refreshInterval: NodeJS.Timeout | null = null;
    
    constructor(private serverManager: ServerManager) {
        // Refresh agents when server state changes
        this.serverManager.onServerStateChanged((isRunning) => {
            if (isRunning) {
                this.startAutoRefresh();
            } else {
                this.stopAutoRefresh();
                this.agents = [];
                this._onDidChangeTreeData.fire();
            }
        });
        
        // Initial load if server is running
        if (this.serverManager.getServerStatus().isRunning) {
            this.startAutoRefresh();
        }
    }
    
    refresh(): void {
        this.loadAgents();
    }
    
    getTreeItem(element: AgentItem | AgentCapabilityItem): vscode.TreeItem {
        return element;
    }
    
    getChildren(element?: AgentItem | AgentCapabilityItem): Thenable<(AgentItem | AgentCapabilityItem)[]> {
        if (!element) {
            // Root level - return agents
            return Promise.resolve(this.getRootItems());
        }
        
        if (element instanceof AgentItem) {
            // Return capabilities for this agent
            return Promise.resolve(
                element.agent.capabilities.map(capability => 
                    new AgentCapabilityItem(capability, vscode.TreeItemCollapsibleState.None)
                )
            );
        }
        
        // Capability items have no children
        return Promise.resolve([]);
    }
    
    private getRootItems(): AgentItem[] {
        if (this.agents.length === 0) {
            return [];
        }
        
        // Sort agents: connected first, then by name
        const sortedAgents = [...this.agents].sort((a, b) => {
            if (a.status !== b.status) {
                return a.status === 'connected' ? -1 : 1;
            }
            return (a.name || a.id).localeCompare(b.name || b.id);
        });
        
        return sortedAgents.map(agent => 
            new AgentItem(
                agent, 
                agent.capabilities.length > 0 
                    ? vscode.TreeItemCollapsibleState.Collapsed 
                    : vscode.TreeItemCollapsibleState.None
            )
        );
    }
    
    private async loadAgents(): Promise<void> {
        try {
            this.agents = await this.serverManager.getAgents();
            this._onDidChangeTreeData.fire();
        } catch (error) {
            console.error('Failed to load agents:', error);
            this.agents = [];
            this._onDidChangeTreeData.fire();
        }
    }
    
    private startAutoRefresh(): void {
        this.stopAutoRefresh();
        
        // Initial load
        this.loadAgents();
        
        // Auto-refresh every 5 seconds
        this.refreshInterval = setInterval(() => {
            this.loadAgents();
        }, 5000);
    }
    
    private stopAutoRefresh(): void {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
        }
    }
    
    // Public methods for external control
    getAgentCount(): number {
        return this.agents.length;
    }
    
    getConnectedAgentCount(): number {
        return this.agents.filter(agent => agent.status === 'connected').length;
    }
    
    getAgentById(id: string): Agent | undefined {
        return this.agents.find(agent => agent.id === id);
    }
    
    dispose(): void {
        this.stopAutoRefresh();
    }
}