import * as vscode from 'vscode';
import { ServerManager } from '../providers/serverManager';

export class DashboardWebviewProvider {
    constructor(
        private readonly _extensionUri: vscode.Uri,
        private readonly serverManager: ServerManager
    ) {}
    
    public setWebviewContent(webview: vscode.Webview) {
        webview.html = this._getHtmlForWebview(webview);
        
        // Handle messages from the webview
        webview.onDidReceiveMessage(async (data) => {
            switch (data.type) {
                case 'getServerStats':
                    await this.sendServerStatsToWebview(webview);
                    break;
                case 'getTools':
                    await this.sendToolsToWebview(webview);
                    break;
                case 'getAgents':
                    await this.sendAgentsToWebview(webview);
                    break;
                case 'startServer':
                    try {
                        await this.serverManager.startServer();
                        vscode.window.showInformationMessage('Server started successfully');
                    } catch (error: any) {
                        vscode.window.showErrorMessage(`Failed to start server: ${error.message}`);
                    }
                    break;
                case 'stopServer':
                    try {
                        await this.serverManager.stopServer();
                        vscode.window.showInformationMessage('Server stopped');
                    } catch (error: any) {
                        vscode.window.showErrorMessage(`Failed to stop server: ${error.message}`);
                    }
                    break;
                case 'restartServer':
                    try {
                        await this.serverManager.restartServer();
                        vscode.window.showInformationMessage('Server restarted');
                    } catch (error: any) {
                        vscode.window.showErrorMessage(`Failed to restart server: ${error.message}`);
                    }
                    break;
            }
        });
        
        // Send initial data
        this.sendInitialDataToWebview(webview);
    }
    
    private async sendInitialDataToWebview(webview: vscode.Webview) {
        await this.sendServerStatsToWebview(webview);
        await this.sendToolsToWebview(webview);
        await this.sendAgentsToWebview(webview);
    }
    
    private async sendServerStatsToWebview(webview: vscode.Webview) {
        try {
            const stats = await this.serverManager.getServerStats();
            const status = this.serverManager.getServerStatus();
            
            webview.postMessage({
                type: 'serverStats',
                stats: stats,
                status: status
            });
        } catch (error) {
            console.error('Failed to send server stats to webview:', error);
        }
    }
    
    private async sendToolsToWebview(webview: vscode.Webview) {
        try {
            const tools = await this.serverManager.getTools();
            webview.postMessage({
                type: 'toolsList',
                tools: tools
            });
        } catch (error) {
            console.error('Failed to send tools to webview:', error);
        }
    }
    
    private async sendAgentsToWebview(webview: vscode.Webview) {
        try {
            const agents = await this.serverManager.getAgents();
            webview.postMessage({
                type: 'agentsList',
                agents: agents
            });
        } catch (error) {
            console.error('Failed to send agents to webview:', error);
        }
    }
    
    private _getHtmlForWebview(webview: vscode.Webview) {
        const styleUri = webview.asWebviewUri(vscode.Uri.joinPath(this._extensionUri, 'media', 'dashboard.css'));
        
        return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link href="${styleUri}" rel="stylesheet">
    <title>AionMCP Dashboard</title>
</head>
<body>
    <div class="dashboard">
        <div class="header">
            <h1>🚀 AionMCP Dashboard</h1>
            <div class="server-controls">
                <button id="startBtn" class="btn btn-success">▶ Start</button>
                <button id="stopBtn" class="btn btn-danger">⏹ Stop</button>
                <button id="restartBtn" class="btn btn-warning">↻ Restart</button>
                <button id="refreshBtn" class="btn btn-secondary">🔄 Refresh</button>
            </div>
        </div>
        
        <div class="status-section">
            <div class="status-card">
                <div class="status-indicator">
                    <div id="statusDot" class="status-dot offline"></div>
                    <span id="statusText">Offline</span>
                </div>
                <div class="status-details">
                    <div>HTTP: <span id="httpPort">-</span></div>
                    <div>gRPC: <span id="grpcPort">-</span></div>
                </div>
            </div>
        </div>
        
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-icon">⏰</div>
                <div class="stat-content">
                    <div class="stat-value" id="uptime">-</div>
                    <div class="stat-label">Uptime</div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-icon">🔧</div>
                <div class="stat-content">
                    <div class="stat-value" id="toolCount">-</div>
                    <div class="stat-label">Tools</div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-icon">👥</div>
                <div class="stat-content">
                    <div class="stat-value" id="agentCount">-</div>
                    <div class="stat-label">Agents</div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-icon">▶</div>
                <div class="stat-content">
                    <div class="stat-value" id="executionCount">-</div>
                    <div class="stat-label">Executions</div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-icon">📊</div>
                <div class="stat-content">
                    <div class="stat-value" id="successRate">-</div>
                    <div class="stat-label">Success Rate</div>
                </div>
            </div>
        </div>
        
        <div class="content-grid">
            <div class="content-section">
                <h3>🔧 Available Tools</h3>
                <div id="toolsList" class="tools-list">
                    <div class="loading">Loading tools...</div>
                </div>
            </div>
            
            <div class="content-section">
                <h3>👥 Connected Agents</h3>
                <div id="agentsList" class="agents-list">
                    <div class="loading">Loading agents...</div>
                </div>
            </div>
        </div>
        
        <div class="logs-section">
            <h3>📋 Recent Activity</h3>
            <div id="activityLogs" class="activity-logs">
                <div class="log-item">
                    <span class="timestamp">--:--:--</span>
                    <span class="message">Dashboard initialized</span>
                </div>
            </div>
        </div>
    </div>
    
    <script>
        const vscode = acquireVsCodeApi();
        
        // Global state
        let isServerRunning = false;
        let autoRefreshInterval = null;
        
        // Initialize
        document.addEventListener('DOMContentLoaded', function() {
            initializeEventListeners();
            requestInitialData();
            startAutoRefresh();
        });
        
        function initializeEventListeners() {
            document.getElementById('startBtn').addEventListener('click', () => {
                vscode.postMessage({ type: 'startServer' });
                addLogMessage('Starting server...');
            });
            
            document.getElementById('stopBtn').addEventListener('click', () => {
                vscode.postMessage({ type: 'stopServer' });
                addLogMessage('Stopping server...');
            });
            
            document.getElementById('restartBtn').addEventListener('click', () => {
                vscode.postMessage({ type: 'restartServer' });
                addLogMessage('Restarting server...');
            });
            
            document.getElementById('refreshBtn').addEventListener('click', requestInitialData);
        }
        
        function requestInitialData() {
            vscode.postMessage({ type: 'getServerStats' });
            vscode.postMessage({ type: 'getTools' });
            vscode.postMessage({ type: 'getAgents' });
        }
        
        function startAutoRefresh() {
            stopAutoRefresh();
            autoRefreshInterval = setInterval(() => {
                if (isServerRunning) {
                    requestInitialData();
                }
            }, 10000); // Refresh every 10 seconds
        }
        
        function stopAutoRefresh() {
            if (autoRefreshInterval) {
                clearInterval(autoRefreshInterval);
                autoRefreshInterval = null;
            }
        }
        
        function updateServerStatus(status, stats) {
            isServerRunning = status.isRunning;
            
            const statusDot = document.getElementById('statusDot');
            const statusText = document.getElementById('statusText');
            const httpPort = document.getElementById('httpPort');
            const grpcPort = document.getElementById('grpcPort');
            
            if (status.isRunning) {
                statusDot.className = 'status-dot online';
                statusText.textContent = 'Online';
                httpPort.textContent = status.port;
                grpcPort.textContent = status.grpcPort;
            } else {
                statusDot.className = 'status-dot offline';
                statusText.textContent = 'Offline';
                httpPort.textContent = '-';
                grpcPort.textContent = '-';
            }
            
            // Update stats if available
            if (stats) {
                document.getElementById('uptime').textContent = formatUptime(stats.uptime);
                document.getElementById('toolCount').textContent = stats.toolCount;
                document.getElementById('agentCount').textContent = stats.agentCount;
                document.getElementById('executionCount').textContent = stats.executionCount;
                document.getElementById('successRate').textContent = \`\${(stats.successRate * 100).toFixed(1)}%\`;
            } else {
                // Clear stats when offline
                document.getElementById('uptime').textContent = '-';
                document.getElementById('toolCount').textContent = '-';
                document.getElementById('agentCount').textContent = '-';
                document.getElementById('executionCount').textContent = '-';
                document.getElementById('successRate').textContent = '-';
            }
            
            // Update button states
            const startBtn = document.getElementById('startBtn');
            const stopBtn = document.getElementById('stopBtn');
            const restartBtn = document.getElementById('restartBtn');
            
            startBtn.disabled = status.isRunning;
            stopBtn.disabled = !status.isRunning;
            restartBtn.disabled = false;
        }
        
        function updateToolsList(tools) {
            const toolsList = document.getElementById('toolsList');
            
            if (!tools || tools.length === 0) {
                toolsList.innerHTML = '<div class="empty-state">No tools available</div>';
                return;
            }
            
            // Group tools by source
            const grouped = groupToolsBySource(tools);
            
            let html = '';
            for (const [source, sourceTools] of Object.entries(grouped)) {
                html += \`<div class="tool-group">
                    <div class="tool-group-header">\${source} (\${sourceTools.length})</div>\`;
                
                sourceTools.forEach(tool => {
                    html += \`<div class="tool-item">
                        <div class="tool-name">\${tool.name}</div>
                        <div class="tool-description">\${tool.description || 'No description'}</div>
                    </div>\`;
                });
                
                html += '</div>';
            }
            
            toolsList.innerHTML = html;
        }
        
        function updateAgentsList(agents) {
            const agentsList = document.getElementById('agentsList');
            
            if (!agents || agents.length === 0) {
                agentsList.innerHTML = '<div class="empty-state">No agents connected</div>';
                return;
            }
            
            let html = '';
            agents.forEach(agent => {
                const statusIcon = agent.status === 'connected' ? '🟢' : '🔴';
                html += \`<div class="agent-item \${agent.status}">
                    <div class="agent-header">
                        <span class="agent-status">\${statusIcon}</span>
                        <span class="agent-name">\${agent.name || agent.id}</span>
                    </div>
                    <div class="agent-capabilities">
                        \${agent.capabilities.map(cap => \`<span class="capability-badge">\${cap}</span>\`).join('')}
                    </div>
                </div>\`;
            });
            
            agentsList.innerHTML = html;
        }
        
        function groupToolsBySource(tools) {
            const grouped = {};
            
            tools.forEach(tool => {
                const source = tool.source || 'Unknown';
                if (!grouped[source]) {
                    grouped[source] = [];
                }
                grouped[source].push(tool);
            });
            
            return grouped;
        }
        
        function formatUptime(seconds) {
            if (!seconds) return '-';
            
            const hours = Math.floor(seconds / 3600);
            const minutes = Math.floor((seconds % 3600) / 60);
            const secs = Math.floor(seconds % 60);
            
            if (hours > 0) {
                return \`\${hours}h \${minutes}m\`;
            } else if (minutes > 0) {
                return \`\${minutes}m \${secs}s\`;
            } else {
                return \`\${secs}s\`;
            }
        }
        
        function addLogMessage(message) {
            const logsContainer = document.getElementById('activityLogs');
            const timestamp = new Date().toLocaleTimeString();
            
            const logItem = document.createElement('div');
            logItem.className = 'log-item';
            logItem.innerHTML = \`
                <span class="timestamp">\${timestamp}</span>
                <span class="message">\${message}</span>
            \`;
            
            logsContainer.insertBefore(logItem, logsContainer.firstChild);
            
            // Keep only last 20 log items
            const logItems = logsContainer.querySelectorAll('.log-item');
            if (logItems.length > 20) {
                logItems[logItems.length - 1].remove();
            }
        }
        
        // Handle messages from extension
        window.addEventListener('message', event => {
            const message = event.data;
            
            switch (message.type) {
                case 'serverStats':
                    updateServerStatus(message.status, message.stats);
                    if (message.stats) {
                        addLogMessage(\`Server stats updated - \${message.stats.executionCount} executions\`);
                    }
                    break;
                case 'toolsList':
                    updateToolsList(message.tools);
                    addLogMessage(\`Loaded \${message.tools.length} tools\`);
                    break;
                case 'agentsList':
                    updateAgentsList(message.agents);
                    const connectedCount = message.agents.filter(a => a.status === 'connected').length;
                    addLogMessage(\`\${connectedCount} agents connected\`);
                    break;
            }
        });
        
        // Cleanup on page unload
        window.addEventListener('beforeunload', stopAutoRefresh);
    </script>
</body>
</html>`;
    }
}