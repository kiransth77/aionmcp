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
exports.DashboardWebviewProvider = void 0;
const vscode = __importStar(require("vscode"));
class DashboardWebviewProvider {
    _extensionUri;
    serverManager;
    static viewType = 'aionmcp.dashboard';
    _view;
    static currentPanel;
    _disposables = [];
    _panel;
    static createOrShow(extensionUri, serverManager) {
        const column = vscode.window.activeTextEditor
            ? vscode.window.activeTextEditor.viewColumn
            : undefined;
        if (DashboardWebviewProvider.currentPanel) {
            DashboardWebviewProvider.currentPanel._panel.reveal(column);
            return;
        }
        const panel = vscode.window.createWebviewPanel(DashboardWebviewProvider.viewType, 'AionMCP Dashboard', column || vscode.ViewColumn.One, {
            enableScripts: true,
            localResourceRoots: [vscode.Uri.joinPath(extensionUri, 'media')],
        });
        DashboardWebviewProvider.currentPanel = new DashboardWebviewProvider(extensionUri, serverManager, panel);
    }
    constructor(_extensionUri, serverManager, panel) {
        this._extensionUri = _extensionUri;
        this.serverManager = serverManager;
        this._panel = panel;
        this._panel.onDidDispose(() => this.dispose(), null, this._disposables);
        this._panel.webview.html = this._getHtmlForWebview(this._panel.webview);
        this.setupWebviewHooks(this._panel.webview);
        this.serverManager.onServerStateChanged((state) => {
            this.sendMessageToWebview({ type: 'serverStatus', status: state ? 'RUNNING' : 'STOPPED' });
        });
        this.sendMessageToWebview({ type: 'serverStatus', status: this.serverManager.getServerStatus().isRunning ? 'RUNNING' : 'STOPPED' });
    }
    resolveWebviewView(webviewView, context, _token) {
        this._view = webviewView;
        webviewView.webview.options = {
            enableScripts: true,
            localResourceRoots: [this._extensionUri]
        };
        webviewView.webview.html = this._getHtmlForWebview(webviewView.webview);
        this.setupWebviewHooks(webviewView.webview);
    }
    setupWebviewHooks(webview) {
        webview.onDidReceiveMessage(async (data) => {
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
        }, null, this._disposables);
        this.sendDashboardData(webview);
    }
    async sendDashboardData(webview) {
        try {
            const data = await this.serverManager.getDashboardData();
            this.sendMessageToWebview({ type: 'dashboardData', data: data }, webview);
        }
        catch (error) {
            this.sendMessageToWebview({ type: 'error', message: `Failed to get dashboard data: ${error.message}` }, webview);
        }
    }
    sendMessageToWebview(message, webview) {
        const targetWebview = webview || this._view?.webview;
        if (targetWebview) {
            targetWebview.postMessage(message);
        }
    }
    dispose() {
        DashboardWebviewProvider.currentPanel = undefined;
        this._disposables.forEach(d => d.dispose());
    }
    _getHtmlForWebview(webview) {
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
exports.DashboardWebviewProvider = DashboardWebviewProvider;
function getNonce() {
    let text = '';
    const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    for (let i = 0; i < 32; i++) {
        text += possible.charAt(Math.floor(Math.random() * possible.length));
    }
    return text;
}
//# sourceMappingURL=dashboardWebview.js.map