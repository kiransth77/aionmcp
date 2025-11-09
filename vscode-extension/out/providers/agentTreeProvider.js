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
exports.AgentTreeProvider = exports.AgentCapabilityItem = exports.AgentItem = void 0;
const vscode = __importStar(require("vscode"));
class AgentItem extends vscode.TreeItem {
    agent;
    collapsibleState;
    constructor(agent, collapsibleState) {
        super(agent.name || agent.id, collapsibleState);
        this.agent = agent;
        this.collapsibleState = collapsibleState;
        this.tooltip = `Agent: ${agent.id}\nStatus: ${agent.status}\nCapabilities: ${agent.capabilities.join(', ')}`;
        this.description = agent.status;
        this.contextValue = 'agent';
        // Set icon and color based on status
        if (agent.status === 'connected') {
            this.iconPath = new vscode.ThemeIcon('account', new vscode.ThemeColor('charts.green'));
        }
        else {
            this.iconPath = new vscode.ThemeIcon('account', new vscode.ThemeColor('charts.red'));
        }
        // Show last seen time for disconnected agents
        if (agent.status === 'disconnected' && agent.lastSeen) {
            const timeDiff = Date.now() - agent.lastSeen.getTime();
            const minutes = Math.floor(timeDiff / 60000);
            const hours = Math.floor(minutes / 60);
            if (hours > 0) {
                this.description = `${agent.status} (${hours}h ago)`;
            }
            else if (minutes > 0) {
                this.description = `${agent.status} (${minutes}m ago)`;
            }
            else {
                this.description = `${agent.status} (just now)`;
            }
        }
    }
}
exports.AgentItem = AgentItem;
class AgentCapabilityItem extends vscode.TreeItem {
    capability;
    collapsibleState;
    constructor(capability, collapsibleState) {
        super(capability, collapsibleState);
        this.capability = capability;
        this.collapsibleState = collapsibleState;
        this.tooltip = `Capability: ${capability}`;
        this.contextValue = 'agentCapability';
        this.iconPath = new vscode.ThemeIcon('symbol-property');
    }
}
exports.AgentCapabilityItem = AgentCapabilityItem;
class AgentTreeProvider {
    serverManager;
    _onDidChangeTreeData = new vscode.EventEmitter();
    onDidChangeTreeData = this._onDidChangeTreeData.event;
    agents = [];
    refreshInterval = null;
    constructor(serverManager) {
        this.serverManager = serverManager;
        // Refresh agents when server state changes
        this.serverManager.onServerStateChanged((isRunning) => {
            if (isRunning) {
                this.startAutoRefresh();
            }
            else {
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
    refresh() {
        this.loadAgents();
    }
    getTreeItem(element) {
        return element;
    }
    getChildren(element) {
        if (!element) {
            // Root level - return agents
            return Promise.resolve(this.getRootItems());
        }
        if (element instanceof AgentItem) {
            // Return capabilities for this agent
            return Promise.resolve(element.agent.capabilities.map(capability => new AgentCapabilityItem(capability, vscode.TreeItemCollapsibleState.None)));
        }
        // Capability items have no children
        return Promise.resolve([]);
    }
    getRootItems() {
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
        return sortedAgents.map(agent => new AgentItem(agent, agent.capabilities.length > 0
            ? vscode.TreeItemCollapsibleState.Collapsed
            : vscode.TreeItemCollapsibleState.None));
    }
    async loadAgents() {
        try {
            this.agents = await this.serverManager.getAgents();
            this._onDidChangeTreeData.fire();
        }
        catch (error) {
            console.error('Failed to load agents:', error);
            this.agents = [];
            this._onDidChangeTreeData.fire();
        }
    }
    startAutoRefresh() {
        this.stopAutoRefresh();
        // Initial load
        this.loadAgents();
        // Auto-refresh every 5 seconds
        this.refreshInterval = setInterval(() => {
            this.loadAgents();
        }, 5000);
    }
    stopAutoRefresh() {
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
            this.refreshInterval = null;
        }
    }
    // Public methods for external control
    getAgentCount() {
        return this.agents.length;
    }
    getConnectedAgentCount() {
        return this.agents.filter(agent => agent.status === 'connected').length;
    }
    getAgentById(id) {
        return this.agents.find(agent => agent.id === id);
    }
    dispose() {
        this.stopAutoRefresh();
    }
}
exports.AgentTreeProvider = AgentTreeProvider;
//# sourceMappingURL=agentTreeProvider.js.map