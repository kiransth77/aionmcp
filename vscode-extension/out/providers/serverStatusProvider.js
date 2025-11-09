"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.ServerStatusProvider = exports.ServerStatusItem = void 0;
const vscode = __importStar(require("vscode"));
class ServerStatusItem extends vscode.TreeItem {
    label;
    value;
    collapsibleState;
    iconName;
    color;
    constructor(label, value, collapsibleState, iconName, color) {
        super(label, collapsibleState);
        this.label = label;
        this.value = value;
        this.collapsibleState = collapsibleState;
        this.iconName = iconName;
        this.color = color;
        this.description = value;
        this.tooltip = `${label}: ${value}`;
        this.contextValue = 'serverStatus';
        if (iconName) {
            this.iconPath = new vscode.ThemeIcon(iconName, color);
        }
    }
}
exports.ServerStatusItem = ServerStatusItem;
class ServerStatusProvider {
    serverManager;
    _onDidChangeTreeData = new vscode.EventEmitter();
    onDidChangeTreeData = this._onDidChangeTreeData.event;
    serverStats = null;
    refreshInterval = null;
    constructor(serverManager) {
        this.serverManager = serverManager;
        // Refresh stats when server state changes
        this.serverManager.onServerStateChanged((isRunning) => {
            if (isRunning) {
                this.startAutoRefresh();
            }
            else {
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
    refresh() {
        this.loadServerStats();
    }
    getTreeItem(element) {
        return element;
    }
    getChildren(element) {
        if (!element) {
            // Root level - return status items
            return Promise.resolve(this.getRootItems());
        }
        // Status items have no children
        return Promise.resolve([]);
    }
    getRootItems() {
        const status = this.serverManager.getServerStatus();
        const items = [];
        // Server running status
        items.push(new ServerStatusItem('Status', status.isRunning ? 'Running' : 'Stopped', vscode.TreeItemCollapsibleState.None, status.isRunning ? 'check' : 'x', status.isRunning
            ? new vscode.ThemeColor('charts.green')
            : new vscode.ThemeColor('charts.red')));
        if (status.isRunning) {
            // Server ports
            items.push(new ServerStatusItem('HTTP Port', status.port.toString(), vscode.TreeItemCollapsibleState.None, 'globe'));
            items.push(new ServerStatusItem('gRPC Port', status.grpcPort.toString(), vscode.TreeItemCollapsibleState.None, 'network'));
            // Server statistics (if available)
            if (this.serverStats) {
                items.push(new ServerStatusItem('Uptime', this.formatUptime(this.serverStats.uptime), vscode.TreeItemCollapsibleState.None, 'clock'));
                items.push(new ServerStatusItem('Tool Count', this.serverStats.toolCount.toString(), vscode.TreeItemCollapsibleState.None, 'tools'));
                items.push(new ServerStatusItem('Connected Agents', this.serverStats.agentCount.toString(), vscode.TreeItemCollapsibleState.None, 'account'));
                items.push(new ServerStatusItem('Executions', this.serverStats.executionCount.toString(), vscode.TreeItemCollapsibleState.None, 'play'));
                items.push(new ServerStatusItem('Success Rate', `${(this.serverStats.successRate * 100).toFixed(1)}%`, vscode.TreeItemCollapsibleState.None, 'graph', this.serverStats.successRate > 0.9
                    ? new vscode.ThemeColor('charts.green')
                    : this.serverStats.successRate > 0.7
                        ? new vscode.ThemeColor('charts.yellow')
                        : new vscode.ThemeColor('charts.red')));
            }
        }
        return items;
    }
    async loadServerStats() {
        try {
            this.serverStats = await this.serverManager.getServerStats();
            this._onDidChangeTreeData.fire();
        }
        catch (error) {
            console.error('Failed to load server stats:', error);
            // Don't clear stats on error, just log it
        }
    }
    formatUptime(uptimeSeconds) {
        const hours = Math.floor(uptimeSeconds / 3600);
        const minutes = Math.floor((uptimeSeconds % 3600) / 60);
        const seconds = Math.floor(uptimeSeconds % 60);
        if (hours > 0) {
            return `${hours}h ${minutes}m ${seconds}s`;
        }
        else if (minutes > 0) {
            return `${minutes}m ${seconds}s`;
        }
        else {
            return `${seconds}s`;
        }
    }
    startAutoRefresh() {
        this.stopAutoRefresh();
        // Initial load
        this.loadServerStats();
        // Auto-refresh every 10 seconds
        this.refreshInterval = setInterval(() => {
            this.loadServerStats();
        }, 10000);
    }
    stopAutoRefresh() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
        }
    }
    dispose() {
        this.stopAutoRefresh();
    }
}
exports.ServerStatusProvider = ServerStatusProvider;
//# sourceMappingURL=serverStatusProvider.js.map