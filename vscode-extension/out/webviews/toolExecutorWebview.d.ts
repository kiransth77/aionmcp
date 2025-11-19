import * as vscode from 'vscode';
import { ServerManager, Tool } from '../providers/serverManager';
export declare class ToolExecutorWebviewProvider implements vscode.WebviewViewProvider {
    private readonly _extensionUri;
    private readonly serverManager;
    private readonly initialTool?;
    static readonly viewType = "aionmcp.toolExecutor";
    private _view?;
    private static currentPanel;
    private readonly _disposables;
    private readonly _panel;
    static createOrShow(extensionUri: vscode.Uri, serverManager: ServerManager, tool?: Tool): void;
    resolveWebviewView(webviewView: vscode.WebviewView, context: vscode.WebviewViewResolveContext, _token: vscode.CancellationToken): void;
    private constructor();
    private updateWithNewTool;
    private setupWebviewHooks;
    private handleToolExecution;
    private sendToolsToWebview;
    private sendToolToWebview;
    dispose(): void;
    private _getHtmlForWebview;
}
//# sourceMappingURL=toolExecutorWebview.d.ts.map