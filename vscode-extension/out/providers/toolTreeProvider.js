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
exports.ToolTreeProvider = exports.ToolCategoryItem = exports.ToolItem = void 0;
const vscode = __importStar(require("vscode"));
class ToolItem extends vscode.TreeItem {
    tool;
    collapsibleState;
    constructor(tool, collapsibleState) {
        super(tool.name, collapsibleState);
        this.tool = tool;
        this.collapsibleState = collapsibleState;
        this.tooltip = tool.description;
        this.description = tool.source || 'unknown';
        this.contextValue = 'tool';
        // Set icon based on tool type
        if (tool.source?.includes('openapi')) {
            this.iconPath = new vscode.ThemeIcon('globe');
        }
        else if (tool.source?.includes('graphql')) {
            this.iconPath = new vscode.ThemeIcon('graph');
        }
        else if (tool.source?.includes('asyncapi')) {
            this.iconPath = new vscode.ThemeIcon('broadcast');
        }
        else {
            this.iconPath = new vscode.ThemeIcon('tools');
        }
        // Add command for double-click execution
        this.command = {
            command: 'aionmcp.openToolExecutor',
            title: 'Execute Tool',
            arguments: [this]
        };
    }
}
exports.ToolItem = ToolItem;
class ToolCategoryItem extends vscode.TreeItem {
    category;
    tools;
    collapsibleState;
    constructor(category, tools, collapsibleState) {
        super(category, collapsibleState);
        this.category = category;
        this.tools = tools;
        this.collapsibleState = collapsibleState;
        this.tooltip = `${tools.length} tools in ${category}`;
        this.description = `${tools.length} tools`;
        this.contextValue = 'toolCategory';
        this.iconPath = new vscode.ThemeIcon('folder');
    }
}
exports.ToolCategoryItem = ToolCategoryItem;
class ToolTreeProvider {
    serverManager;
    _onDidChangeTreeData = new vscode.EventEmitter();
    onDidChangeTreeData = this._onDidChangeTreeData.event;
    tools = [];
    filterText = '';
    constructor(serverManager) {
        this.serverManager = serverManager;
        // Refresh tools when server state changes
        this.serverManager.onServerStateChanged(() => {
            this.refresh();
        });
        // Initial load
        this.refresh();
    }
    refresh() {
        this.loadTools();
    }
    getTreeItem(element) {
        return element;
    }
    getChildren(element) {
        if (!element) {
            // Root level - return categories or tools
            return Promise.resolve(this.getRootItems());
        }
        if (element instanceof ToolCategoryItem) {
            // Return tools in this category
            return Promise.resolve(element.tools.map(tool => new ToolItem(tool, vscode.TreeItemCollapsibleState.None)));
        }
        // Tool items have no children
        return Promise.resolve([]);
    }
    getRootItems() {
        const filteredTools = this.getFilteredTools();
        if (filteredTools.length === 0) {
            return [];
        }
        // Group tools by source/category
        const categories = this.groupToolsByCategory(filteredTools);
        if (categories.size <= 1) {
            // If only one category or no categories, show tools directly
            return filteredTools.map(tool => new ToolItem(tool, vscode.TreeItemCollapsibleState.None));
        }
        // Show categories
        const categoryItems = [];
        for (const [category, tools] of categories) {
            categoryItems.push(new ToolCategoryItem(category, tools, vscode.TreeItemCollapsibleState.Expanded));
        }
        return categoryItems.sort((a, b) => a.category.localeCompare(b.category));
    }
    getFilteredTools() {
        if (!this.filterText) {
            return this.tools;
        }
        const filter = this.filterText.toLowerCase();
        return this.tools.filter(tool => tool.name.toLowerCase().includes(filter) ||
            tool.description.toLowerCase().includes(filter) ||
            (tool.source && tool.source.toLowerCase().includes(filter)));
    }
    groupToolsByCategory(tools) {
        const categories = new Map();
        for (const tool of tools) {
            let category = 'Other';
            if (tool.source) {
                if (tool.source.includes('openapi')) {
                    category = 'OpenAPI';
                }
                else if (tool.source.includes('graphql')) {
                    category = 'GraphQL';
                }
                else if (tool.source.includes('asyncapi')) {
                    category = 'AsyncAPI';
                }
                else {
                    // Use source as category
                    category = tool.source;
                }
            }
            if (!categories.has(category)) {
                categories.set(category, []);
            }
            categories.get(category).push(tool);
        }
        return categories;
    }
    async loadTools() {
        try {
            this.tools = await this.serverManager.getTools();
            this._onDidChangeTreeData.fire();
        }
        catch (error) {
            console.error('Failed to load tools:', error);
            this.tools = [];
            this._onDidChangeTreeData.fire();
        }
    }
    // Public methods for external control
    setFilter(filterText) {
        this.filterText = filterText;
        this._onDidChangeTreeData.fire();
    }
    clearFilter() {
        this.filterText = '';
        this._onDidChangeTreeData.fire();
    }
    getToolCount() {
        return this.tools.length;
    }
    getFilteredToolCount() {
        return this.getFilteredTools().length;
    }
}
exports.ToolTreeProvider = ToolTreeProvider;
//# sourceMappingURL=toolTreeProvider.js.map