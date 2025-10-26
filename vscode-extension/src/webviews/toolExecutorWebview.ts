import * as vscode from 'vscode';
import { ServerManager, Tool } from '../providers/serverManager';

export class ToolExecutorWebviewProvider implements vscode.WebviewViewProvider {
    private _view?: vscode.WebviewView;
    
    constructor(
        private readonly _extensionUri: vscode.Uri,
        private readonly serverManager: ServerManager
    ) {}
    
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
        
        // Handle messages from the webview
        webviewView.webview.onDidReceiveMessage(async (data) => {
            switch (data.type) {
                case 'executeTool':
                    await this.handleToolExecution(data.toolName, data.parameters);
                    break;
                case 'getTools':
                    await this.sendToolsToWebview();
                    break;
                case 'getTool':
                    await this.sendToolToWebview(data.toolName);
                    break;
            }
        });
        
        // Send initial tools list
        this.sendToolsToWebview();
    }
    
    public setWebviewContent(webview: vscode.Webview, tool?: Tool) {
        webview.html = this._getHtmlForWebview(webview, tool);
        
        // Handle messages for standalone webview
        webview.onDidReceiveMessage(async (data) => {
            switch (data.type) {
                case 'executeTool':
                    await this.handleToolExecution(data.toolName, data.parameters);
                    break;
                case 'getTools':
                    await this.sendToolsToWebview(webview);
                    break;
                case 'getTool':
                    await this.sendToolToWebview(data.toolName, webview);
                    break;
            }
        });
        
        if (tool) {
            // Send the specific tool to the webview
            webview.postMessage({
                type: 'toolLoaded',
                tool: tool
            });
        } else {
            // Send tools list
            this.sendToolsToWebview(webview);
        }
    }
    
    private async handleToolExecution(toolName: string, parameters: any) {
        try {
            const result = await this.serverManager.executeTool(toolName, parameters);
            this.sendMessageToWebview({
                type: 'executionResult',
                success: true,
                result: result,
                toolName: toolName
            });
        } catch (error: any) {
            this.sendMessageToWebview({
                type: 'executionResult',
                success: false,
                error: error.message,
                toolName: toolName
            });
        }
    }
    
    private async sendToolsToWebview(webview?: vscode.Webview) {
        try {
            const tools = await this.serverManager.getTools();
            this.sendMessageToWebview({
                type: 'toolsList',
                tools: tools
            }, webview);
        } catch (error) {
            console.error('Failed to send tools to webview:', error);
        }
    }
    
    private async sendToolToWebview(toolName: string, webview?: vscode.Webview) {
        try {
            const tool = await this.serverManager.getTool(toolName);
            if (tool) {
                this.sendMessageToWebview({
                    type: 'toolLoaded',
                    tool: tool
                }, webview);
            }
        } catch (error) {
            console.error('Failed to send tool to webview:', error);
        }
    }
    
    private sendMessageToWebview(message: any, webview?: vscode.Webview) {
        if (webview) {
            webview.postMessage(message);
        } else if (this._view) {
            this._view.webview.postMessage(message);
        }
    }
    
    private _getHtmlForWebview(webview: vscode.Webview, tool?: Tool) {
        const scriptUri = webview.asWebviewUri(vscode.Uri.joinPath(this._extensionUri, 'media', 'toolExecutor.js'));
        const styleUri = webview.asWebviewUri(vscode.Uri.joinPath(this._extensionUri, 'media', 'toolExecutor.css'));
        
        return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link href="${styleUri}" rel="stylesheet">
    <title>Tool Executor</title>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🔧 Tool Executor</h2>
            <button id="refreshTools" class="btn btn-secondary">↻ Refresh</button>
        </div>
        
        <div class="tool-selector">
            <label for="toolSelect">Select Tool:</label>
            <select id="toolSelect" ${tool ? 'disabled' : ''}>
                <option value="">-- Select a tool --</option>
            </select>
        </div>
        
        <div id="toolDetails" class="tool-details" style="display: none;">
            <div class="tool-info">
                <h3 id="toolName"></h3>
                <p id="toolDescription"></p>
                <div class="tool-meta">
                    <span class="badge" id="toolSource"></span>
                    <span class="badge" id="toolVersion"></span>
                </div>
            </div>
            
            <div class="parameters-section">
                <h4>Parameters:</h4>
                <div id="parametersForm"></div>
            </div>
            
            <div class="actions">
                <button id="executeBtn" class="btn btn-primary">▶ Execute Tool</button>
                <button id="clearBtn" class="btn btn-secondary">🗑 Clear</button>
            </div>
        </div>
        
        <div id="executionResults" class="results-section" style="display: none;">
            <h4>Execution Results:</h4>
            <div class="results-header">
                <span id="executionStatus" class="status"></span>
                <span id="executionTime" class="timestamp"></span>
            </div>
            <div id="resultsContent" class="results-content"></div>
        </div>
        
        <div id="executionHistory" class="history-section">
            <h4>Execution History:</h4>
            <div id="historyList" class="history-list">
                <p class="no-history">No executions yet</p>
            </div>
        </div>
    </div>
    
    <script>
        const vscode = acquireVsCodeApi();
        
        // Global state
        let currentTool = null;
        let executionHistory = [];
        
        // Initialize
        document.addEventListener('DOMContentLoaded', function() {
            initializeEventListeners();
            requestTools();
            
            ${tool ? `
            // Pre-load specific tool
            currentTool = ${JSON.stringify(tool)};
            loadTool(currentTool);
            ` : ''}
        });
        
        function initializeEventListeners() {
            document.getElementById('refreshTools').addEventListener('click', requestTools);
            document.getElementById('toolSelect').addEventListener('change', onToolSelected);
            document.getElementById('executeBtn').addEventListener('click', executeTool);
            document.getElementById('clearBtn').addEventListener('click', clearResults);
        }
        
        function requestTools() {
            vscode.postMessage({ type: 'getTools' });
        }
        
        function onToolSelected() {
            const toolSelect = document.getElementById('toolSelect');
            const selectedToolName = toolSelect.value;
            
            if (selectedToolName) {
                vscode.postMessage({ type: 'getTool', toolName: selectedToolName });
            } else {
                hideToolDetails();
            }
        }
        
        function loadTool(tool) {
            currentTool = tool;
            
            document.getElementById('toolName').textContent = tool.name;
            document.getElementById('toolDescription').textContent = tool.description || 'No description available';
            document.getElementById('toolSource').textContent = tool.source || 'Unknown';
            document.getElementById('toolVersion').textContent = tool.version || 'v1.0';
            
            generateParametersForm(tool.inputSchema);
            
            document.getElementById('toolDetails').style.display = 'block';
        }
        
        function hideToolDetails() {
            document.getElementById('toolDetails').style.display = 'none';
            currentTool = null;
        }
        
        function generateParametersForm(schema) {
            const form = document.getElementById('parametersForm');
            form.innerHTML = '';
            
            if (!schema || !schema.properties) {
                form.innerHTML = '<p class="no-params">This tool requires no parameters.</p>';
                return;
            }
            
            const properties = schema.properties;
            const required = schema.required || [];
            
            for (const [propName, propSchema] of Object.entries(properties)) {
                const isRequired = required.includes(propName);
                const fieldDiv = document.createElement('div');
                fieldDiv.className = 'parameter-field';
                
                const label = document.createElement('label');
                label.textContent = propName + (isRequired ? ' *' : '');
                label.setAttribute('for', propName);
                
                const input = createInputElement(propName, propSchema, isRequired);
                
                const description = document.createElement('small');
                description.textContent = propSchema.description || '';
                description.className = 'param-description';
                
                fieldDiv.appendChild(label);
                fieldDiv.appendChild(input);
                if (propSchema.description) {
                    fieldDiv.appendChild(description);
                }
                
                form.appendChild(fieldDiv);
            }
        }
        
        function createInputElement(name, schema, required) {
            let input;
            
            switch (schema.type) {
                case 'boolean':
                    input = document.createElement('input');
                    input.type = 'checkbox';
                    break;
                case 'integer':
                case 'number':
                    input = document.createElement('input');
                    input.type = 'number';
                    if (schema.minimum !== undefined) input.min = schema.minimum;
                    if (schema.maximum !== undefined) input.max = schema.maximum;
                    break;
                case 'array':
                    input = document.createElement('textarea');
                    input.placeholder = 'Enter JSON array, e.g., ["item1", "item2"]';
                    break;
                case 'object':
                    input = document.createElement('textarea');
                    input.placeholder = 'Enter JSON object, e.g., {"key": "value"}';
                    break;
                default:
                    input = document.createElement('input');
                    input.type = 'text';
            }
            
            input.id = name;
            input.name = name;
            input.required = required;
            
            if (schema.default !== undefined) {
                if (schema.type === 'boolean') {
                    input.checked = schema.default;
                } else {
                    input.value = schema.default;
                }
            }
            
            return input;
        }
        
        function executeTool() {
            if (!currentTool) {
                showError('No tool selected');
                return;
            }
            
            const parameters = collectParameters();
            if (parameters === null) {
                return; // Validation failed
            }
            
            showLoading();
            
            vscode.postMessage({
                type: 'executeTool',
                toolName: currentTool.name,
                parameters: parameters
            });
        }
        
        function collectParameters() {
            const form = document.getElementById('parametersForm');
            const inputs = form.querySelectorAll('input, textarea, select');
            const parameters = {};
            
            for (const input of inputs) {
                const name = input.name;
                let value = input.value;
                
                if (input.type === 'checkbox') {
                    value = input.checked;
                } else if (input.type === 'number') {
                    value = value ? Number(value) : undefined;
                } else if (input.tagName === 'TEXTAREA') {
                    if (value.trim()) {
                        try {
                            value = JSON.parse(value);
                        } catch (e) {
                            showError(\`Invalid JSON in parameter "\${name}": \${e.message}\`);
                            return null;
                        }
                    } else {
                        value = undefined;
                    }
                }
                
                if (value !== undefined && value !== '') {
                    parameters[name] = value;
                }
            }
            
            return parameters;
        }
        
        function showLoading() {
            const resultsSection = document.getElementById('executionResults');
            const statusSpan = document.getElementById('executionStatus');
            const timeSpan = document.getElementById('executionTime');
            const contentDiv = document.getElementById('resultsContent');
            
            statusSpan.textContent = '⏳ Executing...';
            statusSpan.className = 'status loading';
            timeSpan.textContent = new Date().toLocaleTimeString();
            contentDiv.innerHTML = '<div class="loading">Executing tool...</div>';
            
            resultsSection.style.display = 'block';
        }
        
        function showExecutionResult(success, result, error, toolName) {
            const statusSpan = document.getElementById('executionStatus');
            const contentDiv = document.getElementById('resultsContent');
            
            if (success) {
                statusSpan.textContent = '✅ Success';
                statusSpan.className = 'status success';
                contentDiv.innerHTML = \`<pre><code>\${JSON.stringify(result, null, 2)}</code></pre>\`;
            } else {
                statusSpan.textContent = '❌ Error';
                statusSpan.className = 'status error';
                contentDiv.innerHTML = \`<div class="error-message">\${error}</div>\`;
            }
            
            // Add to history
            addToHistory(toolName, success, result, error);
        }
        
        function addToHistory(toolName, success, result, error) {
            const historyItem = {
                toolName,
                success,
                result,
                error,
                timestamp: new Date()
            };
            
            executionHistory.unshift(historyItem);
            if (executionHistory.length > 10) {
                executionHistory.pop();
            }
            
            updateHistoryDisplay();
        }
        
        function updateHistoryDisplay() {
            const historyList = document.getElementById('historyList');
            
            if (executionHistory.length === 0) {
                historyList.innerHTML = '<p class="no-history">No executions yet</p>';
                return;
            }
            
            historyList.innerHTML = executionHistory.map(item => \`
                <div class="history-item \${item.success ? 'success' : 'error'}">
                    <div class="history-header">
                        <span class="tool-name">\${item.toolName}</span>
                        <span class="timestamp">\${item.timestamp.toLocaleTimeString()}</span>
                    </div>
                    <div class="history-status">
                        \${item.success ? '✅ Success' : '❌ Error'}
                    </div>
                </div>
            \`).join('');
        }
        
        function clearResults() {
            document.getElementById('executionResults').style.display = 'none';
            document.getElementById('parametersForm').querySelectorAll('input, textarea').forEach(input => {
                if (input.type === 'checkbox') {
                    input.checked = false;
                } else {
                    input.value = '';
                }
            });
        }
        
        function showError(message) {
            const statusSpan = document.getElementById('executionStatus');
            const timeSpan = document.getElementById('executionTime');
            const contentDiv = document.getElementById('resultsContent');
            
            statusSpan.textContent = '❌ Error';
            statusSpan.className = 'status error';
            timeSpan.textContent = new Date().toLocaleTimeString();
            contentDiv.innerHTML = \`<div class="error-message">\${message}</div>\`;
            
            document.getElementById('executionResults').style.display = 'block';
        }
        
        // Handle messages from extension
        window.addEventListener('message', event => {
            const message = event.data;
            
            switch (message.type) {
                case 'toolsList':
                    loadToolsList(message.tools);
                    break;
                case 'toolLoaded':
                    loadTool(message.tool);
                    break;
                case 'executionResult':
                    showExecutionResult(message.success, message.result, message.error, message.toolName);
                    break;
            }
        });
        
        function loadToolsList(tools) {
            const toolSelect = document.getElementById('toolSelect');
            const currentValue = toolSelect.value;
            
            toolSelect.innerHTML = '<option value="">-- Select a tool --</option>';
            
            tools.forEach(tool => {
                const option = document.createElement('option');
                option.value = tool.name;
                option.textContent = \`\${tool.name} (\${tool.source || 'unknown'})\`;
                toolSelect.appendChild(option);
            });
            
            // Restore selection if it still exists
            if (currentValue && tools.some(t => t.name === currentValue)) {
                toolSelect.value = currentValue;
            }
        }
    </script>
</body>
</html>`;
    }
}