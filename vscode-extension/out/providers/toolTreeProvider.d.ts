import * as vscode from 'vscode';
import { ServerManager, Tool } from './serverManager';
export declare class ToolItem extends vscode.TreeItem {
    readonly tool: Tool;
    readonly collapsibleState: vscode.TreeItemCollapsibleState;
    constructor(tool: Tool, collapsibleState: vscode.TreeItemCollapsibleState);
}
export declare class ToolCategoryItem extends vscode.TreeItem {
    readonly category: string;
    readonly tools: Tool[];
    readonly collapsibleState: vscode.TreeItemCollapsibleState;
    constructor(category: string, tools: Tool[], collapsibleState: vscode.TreeItemCollapsibleState);
}
export declare class ToolTreeProvider implements vscode.TreeDataProvider<ToolItem | ToolCategoryItem> {
    private serverManager;
    private _onDidChangeTreeData;
    readonly onDidChangeTreeData: vscode.Event<ToolItem | ToolCategoryItem | undefined | null | void>;
    private tools;
    private filterText;
    constructor(serverManager: ServerManager);
    refresh(): void;
    getTreeItem(element: ToolItem | ToolCategoryItem): vscode.TreeItem;
    getChildren(element?: ToolItem | ToolCategoryItem): Thenable<(ToolItem | ToolCategoryItem)[]>;
    private getRootItems;
    private getFilteredTools;
    private groupToolsByCategory;
    private loadTools;
    setFilter(filterText: string): void;
    clearFilter(): void;
    getToolCount(): number;
    getFilteredToolCount(): number;
}
//# sourceMappingURL=toolTreeProvider.d.ts.map