import * as vscode from 'vscode';
import { ServerManager } from '../providers/serverManager';
export declare class DashboardWebviewProvider implements vscode.WebviewViewProvider {
    private readonly _extensionUri;
    private readonly serverManager;
    static readonly viewType = "aionmcp.dashboard";
    private _view?;
    private static currentPanel;
    private readonly _disposables;
    private readonly _panel;
    static createOrShow(extensionUri: vscode.Uri, serverManager: ServerManager): void;
    private constructor();
    resolveWebviewView(webviewView: vscode.WebviewView, context: vscode.WebviewViewResolveContext, _token: vscode.CancellationToken): void;
    private setupWebviewHooks;
    private sendDashboardData;
    private sendMessageToWebview;
    dispose(): void;
    private _getHtmlForWebview;
}
//# sourceMappingURL=dashboardWebview.d.ts.map