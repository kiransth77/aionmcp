import * as vscode from 'vscode';
import { ServerManager } from '../providers/serverManager';

export class DashboardWebviewProvider implements vscode.WebviewViewProvider {
    public static readonly viewType = 'aionmcp.dashboard';

    private _view?: vscode.WebviewView;
    private static currentPanel: DashboardWebviewProvider | undefined;
    private readonly _disposables: vscode.Disposable[] = [];
    private readonly _panel: vscode.WebviewPanel;

    public static createOrShow(extensionUri: vscode.Uri, serverManager: ServerManager) {
        const column = vscode.window.activeTextEditor
            ? vscode.window.activeTextEditor.viewColumn
            : undefined;

        if (DashboardWebviewProvider.currentPanel) {
            DashboardWebviewProvider.currentPanel._panel.reveal(column);
            return;
        }

        const panel = vscode.window.createWebviewPanel(
            DashboardWebviewProvider.viewType,
            'AionMCP Dashboard',
            column || vscode.ViewColumn.One,
            {
                enableScripts: true,
                localResourceRoots: [vscode.Uri.joinPath(extensionUri, 'media')],
            }
        );

        DashboardWebviewProvider.currentPanel = new DashboardWebviewProvider(extensionUri, serverManager, panel);
    }

    private constructor(
        private readonly _extensionUri: vscode.Uri,
        private readonly serverManager: ServerManager,
        panel: vscode.WebviewPanel
    ) {
        this._panel = panel;
        this._panel.onDidDispose(() => this.dispose(), null, this._disposables);
        this._panel.webview.html = this._getHtmlForWebview(this._panel.webview);
        this.setupWebviewHooks(this._panel.webview);

        this.serverManager.onServerStateChanged((state: boolean) => {
            this.sendMessageToWebview({ type: 'serverStatus', status: state ? 'RUNNING' : 'STOPPED' });
        });
        this.sendMessageToWebview({ type: 'serverStatus', status: this.serverManager.getServerStatus().isRunning ? 'RUNNING' : 'STOPPED' });
    }

    public resolveWebviewView(
        webviewView: vscode.WebviewView,
        context: vscode.WebviewViewResolveContext,
        _token: vscode.CancellationToken,
    ) {
        this._view = webviewView;

        webviewView.webview.options = {
            enableScripts: true,
            localResourceRoots: [this._extensionUri]
        };

        webviewView.webview.html = this._getHtmlForWebview(webviewView.webview);
        this.setupWebviewHooks(webviewView.webview);
    }

    private setupWebviewHooks(webview: vscode.Webview) {
        webview.onDidReceiveMessage(
            async (data) => {
                switch (data.type) {
                    case 'getDashboardData':
                        await this.sendDashboardData(webview);
                        break;
                    case 'startServer':
                        await this.serverManager.startServer();
                        break;
                    case 'stopServer':
                        await this.serverManager.stopServer();
                        break;
                    case 'restartServer':
                        await this.serverManager.restartServer();
                        break;
                }
            },
            null,
            this._disposables
        );

        this.sendDashboardData(webview);
    }

    private async sendDashboardData(webview: vscode.Webview) {
        try {
            const data = await this.serverManager.getDashboardData();
            this.sendMessageToWebview({ type: 'dashboardData', data: data }, webview);
        } catch (error: any) {
            this.sendMessageToWebview({ type: 'error', message: `Failed to get dashboard data: ${error.message}` }, webview);
        }
    }

    private sendMessageToWebview(message: any, webview?: vscode.Webview) {
        const targetWebview = webview || this._view?.webview;
        if (targetWebview) {
            targetWebview.postMessage(message);
        }
    }

    public dispose() {
        DashboardWebviewProvider.currentPanel = undefined;
        this._panel.dispose();
        this._disposables.forEach(d => d.dispose());
    }

    private _getHtmlForWebview(webview: vscode.Webview): string {
        const styleUri = webview.asWebviewUri(vscode.Uri.joinPath(this._extensionUri, 'media', 'dashboard.css'));
        const nonce = getNonce();

        return `<!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <title>AionMCP Dashboard</title>
                <link href="${styleUri}" rel="stylesheet">
            </head>
            <body>
                <div class="container">
                    <h1>AionMCP Dashboard</h1>
                    <div class="server-status">
                        Server Status: <span id="status-indicator" class="stopped"></span> <span id="status-text">Not Running</span>
                    </div>
                    <div class="server-controls">
                        <button id="start-btn">Start Server</button>
                        <button id="stop-btn">Stop Server</button>
                        <button id="restart-btn">Restart Server</button>
                    </div>
                    <div id="dashboard-data"></div>
                </div>
                <script nonce="${nonce}">
                    const vscode = acquireVsCodeApi();

                    const statusIndicator = document.getElementById('status-indicator');
                    const statusText = document.getElementById('status-text');
                    const startBtn = document.getElementById('start-btn');
                    const stopBtn = document.getElementById('stop-btn');
                    const restartBtn = document.getElementById('restart-btn');
                    const dashboardDataEl = document.getElementById('dashboard-data');

                    startBtn.addEventListener('click', () => vscode.postMessage({ type: 'startServer' }));
                    stopBtn.addEventListener('click', () => vscode.postMessage({ type: 'stopServer' }));
                    restartBtn.addEventListener('click', () => vscode.postMessage({ type: 'restartServer' }));

                    window.addEventListener('message', event => {
                        const message = event.data;
                        switch (message.type) {
                            case 'serverStatus':
                                updateServerStatus(message.status);
                                break;
                            case 'dashboardData':
                                renderDashboardData(message.data);
                                break;
                            case 'error':
                                dashboardDataEl.innerHTML = \`<p class="error">\${message.message}</p>\`;
                                break;
                        }
                    });

                    function updateServerStatus(status) {
                        const isRunning = status === 'RUNNING';
                        statusIndicator.className = isRunning ? 'running' : 'stopped';
                        statusText.textContent = isRunning ? 'Running' : 'Not Running';
                        startBtn.disabled = isRunning;
                        stopBtn.disabled = !isRunning;
                        restartBtn.disabled = !isRunning;
                    }

                    function renderDashboardData(data) {
                        if (!data) {
                            dashboardDataEl.innerHTML = '<p>No data to display.</p>';
                            return;
                        }
                        let content = '<h2>Server Info</h2>';
                        content += \`<ul>
                            <li><strong>Uptime:</strong> \${data.uptime}s</li>
                            <li><strong>Tools Registered:</strong> \${data.toolCount}</li>
                            <li><strong>Agents Connected:</strong> \${data.agentCount}</li>
                            <li><strong>Total Executions:</strong> \${data.executionCount}</li>
                            <li><strong>Success Rate:</strong> \${(data.successRate * 100).toFixed(2)}%</li>
                        </ul>\`;
                        dashboardDataEl.innerHTML = content;
                    }

                    // Request initial data
                    vscode.postMessage({ type: 'getDashboardData' });
                </script>
            </body>
            </html>`;
    }
}

function getNonce() {
    let text = '';
    const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    for (let i = 0; i < 32; i++) {
        text += possible.charAt(Math.floor(Math.random() * possible.length));
    }
    return text;
}