import * as vscode from 'vscode';
import { ServerManager, Tool } from '../providers/serverManager';
export declare class ToolExecutorWebviewProvider implements vscode.WebviewViewProvider {
    private readonly _extensionUri;
    private readonly serverManager;
    private _view?;
    constructor(_extensionUri: vscode.Uri, serverManager: ServerManager);
    resolveWebviewView(webviewView: vscode.WebviewView, context: vscode.WebviewViewResolveContext, _token: vscode.CancellationToken): void;
    setWebviewContent(webview: vscode.Webview, tool?: Tool): void;
    private handleToolExecution;
    private sendToolsToWebview;
    private sendToolToWebview;
    private sendMessageToWebview;
    private _getHtmlForWebview;
}
//# sourceMappingURL=toolExecutorWebview.d.ts.map