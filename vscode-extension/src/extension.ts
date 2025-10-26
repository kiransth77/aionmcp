import * as vscode from 'vscode';
import { ServerManager } from './providers/serverManager';
import { ToolTreeProvider, ToolItem } from './providers/toolTreeProvider';
import { AgentTreeProvider } from './providers/agentTreeProvider';
import { ServerStatusProvider } from './providers/serverStatusProvider';
import { ToolExecutorWebviewProvider } from './webviews/toolExecutorWebview';
import { LogOutputProvider } from './providers/logOutputProvider';
import { DashboardWebviewProvider } from './webviews/dashboardWebview';

let serverManager: ServerManager;
let toolTreeProvider: ToolTreeProvider;
let agentTreeProvider: AgentTreeProvider;
let serverStatusProvider: ServerStatusProvider;
let logOutputProvider: LogOutputProvider;

export function activate(context: vscode.ExtensionContext) {
    console.log('AionMCP extension is now active!');
    
    // Initialize providers
    serverManager = new ServerManager(context);
    toolTreeProvider = new ToolTreeProvider(serverManager);
    agentTreeProvider = new AgentTreeProvider(serverManager);
    serverStatusProvider = new ServerStatusProvider(serverManager);
    logOutputProvider = new LogOutputProvider();
    
    // Register tree data providers
    vscode.window.registerTreeDataProvider('aionmcp.toolsView', toolTreeProvider);
    vscode.window.registerTreeDataProvider('aionmcp.agentsView', agentTreeProvider);
    vscode.window.registerTreeDataProvider('aionmcp.serverView', serverStatusProvider);
    
    // Register webview providers
    const toolExecutorProvider = new ToolExecutorWebviewProvider(context.extensionUri, serverManager);
    const dashboardProvider = new DashboardWebviewProvider(context.extensionUri, serverManager);
    
    context.subscriptions.push(
        vscode.window.registerWebviewViewProvider(
            'aionmcp.toolExecutor',
            toolExecutorProvider
        )
    );
    
    // Register commands
    registerCommands(context);
    
    // Set extension as active
    vscode.commands.executeCommand('setContext', 'aionmcp.extensionActive', true);
    
    // Auto-start server if configured
    const config = vscode.workspace.getConfiguration('aionmcp');
    if (config.get<boolean>('autoStart', false)) {
        vscode.commands.executeCommand('aionmcp.startServer');
    }
    
    // Set up status bar
    const statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100);
    statusBarItem.text = '$(server-process) AionMCP';
    statusBarItem.tooltip = 'AionMCP Server Status';
    statusBarItem.command = 'aionmcp.showDashboard';
    statusBarItem.show();
    context.subscriptions.push(statusBarItem);
    
    // Update status bar based on server state
    serverManager.onServerStateChanged((isRunning: boolean) => {
        vscode.commands.executeCommand('setContext', 'aionmcp.serverRunning', isRunning);
        statusBarItem.text = isRunning 
            ? '$(server-process) AionMCP $(check)' 
            : '$(server-process) AionMCP $(x)';
        statusBarItem.color = isRunning ? undefined : new vscode.ThemeColor('statusBarItem.errorForeground');
    });
}

function registerCommands(context: vscode.ExtensionContext) {
    // Server management commands
    context.subscriptions.push(
        vscode.commands.registerCommand('aionmcp.startServer', async () => {
            try {
                await serverManager.startServer();
                vscode.window.showInformationMessage('AionMCP server started successfully');
            } catch (error) {
                vscode.window.showErrorMessage(`Failed to start server: ${error}`);
            }
        })
    );
    
    context.subscriptions.push(
        vscode.commands.registerCommand('aionmcp.stopServer', async () => {
            try {
                await serverManager.stopServer();
                vscode.window.showInformationMessage('AionMCP server stopped');
            } catch (error) {
                vscode.window.showErrorMessage(`Failed to stop server: ${error}`);
            }
        })
    );
    
    context.subscriptions.push(
        vscode.commands.registerCommand('aionmcp.restartServer', async () => {
            try {
                await serverManager.restartServer();
                vscode.window.showInformationMessage('AionMCP server restarted');
            } catch (error) {
                vscode.window.showErrorMessage(`Failed to restart server: ${error}`);
            }
        })
    );
    
    // Tool management commands
    context.subscriptions.push(
        vscode.commands.registerCommand('aionmcp.refreshTools', () => {
            toolTreeProvider.refresh();
        })
    );
    
    context.subscriptions.push(
        vscode.commands.registerCommand('aionmcp.executeTool', async (toolItem: ToolItem) => {
            if (!toolItem) {
                vscode.window.showErrorMessage('No tool selected');
                return;
            }
            
            try {
                // Open quick tool execution
                const result = await executeToolQuick(toolItem);
                if (result) {
                    vscode.window.showInformationMessage(`Tool executed successfully: ${toolItem.tool.name}`);
                }
            } catch (error) {
                vscode.window.showErrorMessage(`Tool execution failed: ${error}`);
            }
        })
    );
    
    context.subscriptions.push(
        vscode.commands.registerCommand('aionmcp.openToolExecutor', async (toolItem: ToolItem) => {
            // Open detailed tool executor webview
            const panel = vscode.window.createWebviewPanel(
                'aionmcp.toolExecutor',
                `Execute: ${toolItem ? toolItem.tool.name : 'Tool Executor'}`,
                vscode.ViewColumn.One,
                {
                    enableScripts: true,
                    retainContextWhenHidden: true
                }
            );
            
            const provider = new ToolExecutorWebviewProvider(context.extensionUri, serverManager);
            provider.setWebviewContent(panel.webview, toolItem?.tool);
        })
    );
    
    // API Spec import
    context.subscriptions.push(
        vscode.commands.registerCommand('aionmcp.importApiSpec', async () => {
            const fileUri = await vscode.window.showOpenDialog({
                canSelectFiles: true,
                canSelectFolders: false,
                canSelectMany: false,
                filters: {
                    'API Specifications': ['json', 'yaml', 'yml', 'graphql']
                },
                openLabel: 'Import API Specification'
            });
            
            if (fileUri && fileUri[0]) {
                try {
                    await serverManager.importApiSpec(fileUri[0].fsPath);
                    vscode.window.showInformationMessage('API specification imported successfully');
                    toolTreeProvider.refresh();
                } catch (error) {
                    vscode.window.showErrorMessage(`Failed to import API spec: ${error}`);
                }
            }
        })
    );
    
    // Logging and monitoring
    context.subscriptions.push(
        vscode.commands.registerCommand('aionmcp.viewLogs', () => {
            logOutputProvider.show();
        })
    );
    
    context.subscriptions.push(
        vscode.commands.registerCommand('aionmcp.showDashboard', () => {
            const panel = vscode.window.createWebviewPanel(
                'aionmcp.dashboard',
                'AionMCP Dashboard',
                vscode.ViewColumn.One,
                {
                    enableScripts: true,
                    retainContextWhenHidden: true
                }
            );
            
            const provider = new DashboardWebviewProvider(context.extensionUri, serverManager);
            provider.setWebviewContent(panel.webview);
        })
    );
}

async function executeToolQuick(toolItem: ToolItem): Promise<boolean> {
    // Simple tool execution with basic parameter input
    const tool = toolItem.tool;
    
    // If tool has no parameters, execute directly
    if (!tool.inputSchema || Object.keys(tool.inputSchema.properties || {}).length === 0) {
        const result = await serverManager.executeTool(tool.name, {});
        showExecutionResult(tool.name, result);
        return true;
    }
    
    // For tools with parameters, show quick input
    const params: any = {};
    const properties = tool.inputSchema.properties || {};
    
    for (const [propName, propSchema] of Object.entries(properties)) {
        const schema = propSchema as any;
        const value = await vscode.window.showInputBox({
            prompt: `Enter value for ${propName}`,
            placeHolder: schema.description || `Value for ${propName}`,
            value: schema.default ? String(schema.default) : undefined
        });
        
        if (value === undefined) {
            return false; // User cancelled
        }
        
        // Simple type conversion
        if (schema.type === 'number' || schema.type === 'integer') {
            params[propName] = Number(value);
        } else if (schema.type === 'boolean') {
            params[propName] = value.toLowerCase() === 'true';
        } else {
            params[propName] = value;
        }
    }
    
    const result = await serverManager.executeTool(tool.name, params);
    showExecutionResult(tool.name, result);
    return true;
}

function showExecutionResult(toolName: string, result: any) {
    const resultStr = typeof result === 'string' ? result : JSON.stringify(result, null, 2);
    
    // Show result in output channel
    logOutputProvider.appendLine(`\n=== Tool Execution: ${toolName} ===`);
    logOutputProvider.appendLine(resultStr);
    logOutputProvider.appendLine('=== End Result ===\n');
    logOutputProvider.show();
}

export function deactivate() {
    if (serverManager) {
        serverManager.dispose();
    }
}