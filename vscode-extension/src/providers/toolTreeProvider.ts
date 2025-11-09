import * as vscode from 'vscode';
import { ServerManager, Tool } from './serverManager';

export class ToolItem extends vscode.TreeItem {
    constructor(
        public readonly tool: Tool,
        public readonly collapsibleState: vscode.TreeItemCollapsibleState
    ) {
        super(tool.name, collapsibleState);
        
        this.tooltip = tool.description;
        this.description = tool.source || 'unknown';
        this.contextValue = 'tool';
        
        // Set icon based on tool type
        if (tool.source?.includes('openapi')) {
            this.iconPath = new vscode.ThemeIcon('globe');
        } else if (tool.source?.includes('graphql')) {
            this.iconPath = new vscode.ThemeIcon('graph');
        } else if (tool.source?.includes('asyncapi')) {
            this.iconPath = new vscode.ThemeIcon('broadcast');
        } else {
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

export class ToolCategoryItem extends vscode.TreeItem {
    constructor(
        public readonly category: string,
        public readonly tools: Tool[],
        public readonly collapsibleState: vscode.TreeItemCollapsibleState
    ) {
        super(category, collapsibleState);
        
        this.tooltip = `${tools.length} tools in ${category}`;
        this.description = `${tools.length} tools`;
        this.contextValue = 'toolCategory';
        this.iconPath = new vscode.ThemeIcon('folder');
    }
}

export class ToolTreeProvider implements vscode.TreeDataProvider<ToolItem | ToolCategoryItem> {
    private _onDidChangeTreeData: vscode.EventEmitter<ToolItem | ToolCategoryItem | undefined | null | void> = new vscode.EventEmitter<ToolItem | ToolCategoryItem | undefined | null | void>();
    readonly onDidChangeTreeData: vscode.Event<ToolItem | ToolCategoryItem | undefined | null | void> = this._onDidChangeTreeData.event;
    
    private tools: Tool[] = [];
    private filterText = '';
    
    constructor(private serverManager: ServerManager) {
        // Refresh tools when server state changes
        this.serverManager.onServerStateChanged(() => {
            this.refresh();
        });
        
        // Initial load
        this.refresh();
    }
    
    refresh(): void {
        this.loadTools();
    }
    
    getTreeItem(element: ToolItem | ToolCategoryItem): vscode.TreeItem {
        return element;
    }
    
    getChildren(element?: ToolItem | ToolCategoryItem): Thenable<(ToolItem | ToolCategoryItem)[]> {
        if (!element) {
            // Root level - return categories or tools
            return Promise.resolve(this.getRootItems());
        }
        
        if (element instanceof ToolCategoryItem) {
            // Return tools in this category
            return Promise.resolve(
                element.tools.map(tool => new ToolItem(tool, vscode.TreeItemCollapsibleState.None))
            );
        }
        
        // Tool items have no children
        return Promise.resolve([]);
    }
    
    private getRootItems(): (ToolItem | ToolCategoryItem)[] {
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
        const categoryItems: ToolCategoryItem[] = [];
        for (const [category, tools] of categories) {
            categoryItems.push(new ToolCategoryItem(
                category,
                tools,
                vscode.TreeItemCollapsibleState.Expanded
            ));
        }
        
        return categoryItems.sort((a, b) => a.category.localeCompare(b.category));
    }
    
    private getFilteredTools(): Tool[] {
        if (!this.filterText) {
            return this.tools;
        }
        
        const filter = this.filterText.toLowerCase();
        return this.tools.filter(tool =>
            tool.name.toLowerCase().includes(filter) ||
            tool.description.toLowerCase().includes(filter) ||
            (tool.source && tool.source.toLowerCase().includes(filter))
        );
    }
    
    private groupToolsByCategory(tools: Tool[]): Map<string, Tool[]> {
        const categories = new Map<string, Tool[]>();
        
        for (const tool of tools) {
            let category = 'Other';
            
            if (tool.source) {
                if (tool.source.includes('openapi')) {
                    category = 'OpenAPI';
                } else if (tool.source.includes('graphql')) {
                    category = 'GraphQL';
                } else if (tool.source.includes('asyncapi')) {
                    category = 'AsyncAPI';
                } else {
                    // Use source as category
                    category = tool.source;
                }
            }
            
            if (!categories.has(category)) {
                categories.set(category, []);
            }
            categories.get(category)!.push(tool);
        }
        
        return categories;
    }
    
    private async loadTools(): Promise<void> {
        try {
            this.tools = await this.serverManager.getTools();
            this._onDidChangeTreeData.fire();
        } catch (error) {
            console.error('Failed to load tools:', error);
            this.tools = [];
            this._onDidChangeTreeData.fire();
        }
    }
    
    // Public methods for external control
    setFilter(filterText: string): void {
        this.filterText = filterText;
        this._onDidChangeTreeData.fire();
    }
    
    clearFilter(): void {
        this.filterText = '';
        this._onDidChangeTreeData.fire();
    }
    
    getToolCount(): number {
        return this.tools.length;
    }
    
    getFilteredToolCount(): number {
        return this.getFilteredTools().length;
    }
}