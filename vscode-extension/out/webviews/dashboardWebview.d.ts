import * as vscode from 'vscode';
import { ServerManager } from '../providers/serverManager';
export declare class DashboardWebviewProvider {
    private readonly _extensionUri;
    private readonly serverManager;
    constructor(_extensionUri: vscode.Uri, serverManager: ServerManager);
    setWebviewContent(webview: vscode.Webview): void;
    private sendInitialDataToWebview;
    private sendServerStatsToWebview;
    private sendToolsToWebview;
    private sendAgentsToWebview;
    private _getHtmlForWebview;
}
//# sourceMappingURL=dashboardWebview.d.ts.map