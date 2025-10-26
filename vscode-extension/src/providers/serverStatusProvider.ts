import * as vscode from 'vscode';
import { ServerManager, ServerStats } from './serverManager';

export class ServerStatusItem extends vscode.TreeItem {
    constructor(
        public readonly label: string,
        public readonly value: string,
        public readonly collapsibleState: vscode.TreeItemCollapsibleState,
        public readonly iconName?: string,
        public readonly color?: vscode.ThemeColor
    ) {
        super(label, collapsibleState);
        
        this.description = value;
        this.tooltip = `${label}: ${value}`;
        this.contextValue = 'serverStatus';
        
        if (iconName) {
            this.iconPath = new vscode.ThemeIcon(iconName, color);
        }
    }
}

export class ServerStatusProvider implements vscode.TreeDataProvider<ServerStatusItem> {
    private _onDidChangeTreeData: vscode.EventEmitter<ServerStatusItem | undefined | null | void> = new vscode.EventEmitter<ServerStatusItem | undefined | null | void>();
    readonly onDidChangeTreeData: vscode.Event<ServerStatusItem | undefined | null | void> = this._onDidChangeTreeData.event;
    
    private serverStats: ServerStats | null = null;
    private refreshInterval: NodeJS.Timeout | null = null;
    
    constructor(private serverManager: ServerManager) {
        // Refresh stats when server state changes
        this.serverManager.onServerStateChanged((isRunning) => {
            if (isRunning) {
                this.startAutoRefresh();
            } else {
                this.stopAutoRefresh();
                this.serverStats = null;
                this._onDidChangeTreeData.fire();
            }
        });
        
        // Initial load if server is running
        if (this.serverManager.getServerStatus().isRunning) {
            this.startAutoRefresh();
        }
    }
    
    refresh(): void {
        this.loadServerStats();
    }
    
    getTreeItem(element: ServerStatusItem): vscode.TreeItem {
        return element;
    }
    
    getChildren(element?: ServerStatusItem): Thenable<ServerStatusItem[]> {
        if (!element) {
            // Root level - return status items
            return Promise.resolve(this.getRootItems());
        }
        
        // Status items have no children
        return Promise.resolve([]);
    }
    
    private getRootItems(): ServerStatusItem[] {
        const status = this.serverManager.getServerStatus();
        const items: ServerStatusItem[] = [];
        
        // Server running status
        items.push(new ServerStatusItem(
            'Status',
            status.isRunning ? 'Running' : 'Stopped',
            vscode.TreeItemCollapsibleState.None,
            status.isRunning ? 'check' : 'x',
            status.isRunning 
                ? new vscode.ThemeColor('charts.green')
                : new vscode.ThemeColor('charts.red')
        ));
        
        if (status.isRunning) {
            // Server ports
            items.push(new ServerStatusItem(
                'HTTP Port',
                status.port.toString(),
                vscode.TreeItemCollapsibleState.None,
                'globe'
            ));
            
            items.push(new ServerStatusItem(
                'gRPC Port',
                status.grpcPort.toString(),
                vscode.TreeItemCollapsibleState.None,
                'network'
            ));
            
            // Server statistics (if available)
            if (this.serverStats) {
                items.push(new ServerStatusItem(
                    'Uptime',
                    this.formatUptime(this.serverStats.uptime),
                    vscode.TreeItemCollapsibleState.None,
                    'clock'
                ));
                
                items.push(new ServerStatusItem(
                    'Tool Count',
                    this.serverStats.toolCount.toString(),
                    vscode.TreeItemCollapsibleState.None,
                    'tools'
                ));
                
                items.push(new ServerStatusItem(
                    'Connected Agents',
                    this.serverStats.agentCount.toString(),
                    vscode.TreeItemCollapsibleState.None,
                    'account'
                ));
                
                items.push(new ServerStatusItem(
                    'Executions',
                    this.serverStats.executionCount.toString(),
                    vscode.TreeItemCollapsibleState.None,
                    'play'
                ));
                
                items.push(new ServerStatusItem(
                    'Success Rate',
                    `${(this.serverStats.successRate * 100).toFixed(1)}%`,
                    vscode.TreeItemCollapsibleState.None,
                    'graph',
                    this.serverStats.successRate > 0.9 
                        ? new vscode.ThemeColor('charts.green')
                        : this.serverStats.successRate > 0.7
                            ? new vscode.ThemeColor('charts.yellow')
                            : new vscode.ThemeColor('charts.red')
                ));
            }
        }
        
        return items;
    }
    
    private async loadServerStats(): Promise<void> {
        try {
            this.serverStats = await this.serverManager.getServerStats();
            this._onDidChangeTreeData.fire();
        } catch (error) {
            console.error('Failed to load server stats:', error);
            // Don't clear stats on error, just log it
        }
    }
    
    private formatUptime(uptimeSeconds: number): string {
        const hours = Math.floor(uptimeSeconds / 3600);
        const minutes = Math.floor((uptimeSeconds % 3600) / 60);
        const seconds = Math.floor(uptimeSeconds % 60);
        
        if (hours > 0) {
            return `${hours}h ${minutes}m ${seconds}s`;
        } else if (minutes > 0) {
            return `${minutes}m ${seconds}s`;
        } else {
            return `${seconds}s`;
        }
    }
    
    private startAutoRefresh(): void {
        this.stopAutoRefresh();
        
        // Initial load
        this.loadServerStats();
        
        // Auto-refresh every 10 seconds
        this.refreshInterval = setInterval(() => {
            this.loadServerStats();
        }, 10000);
    }
    
    private stopAutoRefresh(): void {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
        }
    }
    
    dispose(): void {
        this.stopAutoRefresh();
    }
}