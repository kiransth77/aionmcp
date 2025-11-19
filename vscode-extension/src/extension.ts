import * as vscode from 'vscode';
import { ServerManager } from './providers/serverManager';
import { ToolTreeProvider, ToolItem } from './providers/toolTreeProvider';
import { AgentTreeProvider } from './providers/agentTreeProvider';
import { ServerStatusProvider } from './providers/serverStatusProvider';
import { LogOutputProvider } from './providers/logOutputProvider';
import { DashboardWebviewProvider } from './webviews/dashboardWebview';
import { ToolExecutorWebviewProvider } from './webviews/toolExecutorWebview';

let serverManager: ServerManager;
let toolTreeProvider: ToolTreeProvider;
let agentTreeProvider: AgentTreeProvider;
let serverStatusProvider: ServerStatusProvider;
let logOutputProvider: LogOutputProvider;

export function activate(context: vscode.ExtensionContext) {
    console.log('AionMCP extension activating...');
    
    try {
        // Initialize server manager first
        serverManager = new ServerManager(context);
        console.log('ServerManager initialized');
        
        // Initialize tree providers
        toolTreeProvider = new ToolTreeProvider(serverManager);
        console.log('ToolTreeProvider initialized');
        
        agentTreeProvider = new AgentTreeProvider(serverManager);
        console.log('AgentTreeProvider initialized');
        
        serverStatusProvider = new ServerStatusProvider(serverManager);
        console.log('ServerStatusProvider initialized');
        
        logOutputProvider = new LogOutputProvider();
        console.log('LogOutputProvider initialized');
        
        // Register tree data providers
        vscode.window.registerTreeDataProvider('aionmcp.toolsView', toolTreeProvider);
        vscode.window.registerTreeDataProvider('aionmcp.agentsView', agentTreeProvider);
        vscode.window.registerTreeDataProvider('aionmcp.serverView', serverStatusProvider);
        console.log('Tree data providers registered');
        
        // Register commands
        registerCommands(context);
        console.log('Commands registered');
        
        // Set extension as active
        vscode.commands.executeCommand('setContext', 'aionmcp.extensionActive', true);
        
        // Set up status bar
        const statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100);
        statusBarItem.text = '$(server-process) AionMCP';
        statusBarItem.tooltip = 'AionMCP Server Status';
        statusBarItem.command = 'aionmcp.viewLogs';
        statusBarItem.show();
        context.subscriptions.push(statusBarItem);
        console.log('Status bar created');
        
        // Update status bar based on server state
        serverManager.onServerStateChanged((isRunning: boolean) => {
            vscode.commands.executeCommand('setContext', 'aionmcp.serverRunning', isRunning);
            statusBarItem.text = isRunning 
                ? '$(server-process) AionMCP $(check)' 
                : '$(server-process) AionMCP $(x)';
            statusBarItem.color = isRunning ? undefined : new vscode.ThemeColor('statusBarItem.errorForeground');
        });
        
        console.log('AionMCP extension activated successfully');
        vscode.window.showInformationMessage('AionMCP extension is ready!');
        
    } catch (error) {
        const msg = error instanceof Error ? error.message : String(error);
        console.error('Extension activation error:', msg, error);
        vscode.window.showErrorMessage(`Extension failed: ${msg}`);
    }
}

function registerCommands(context: vscode.ExtensionContext) {
    const commands: Record<string, any> = {
        'aionmcp.startServer': async () => {
            console.log('Start command received');
            try {
                await serverManager.startServer();
            } catch (error) {
                const msg = error instanceof Error ? error.message : String(error);
                console.error('Start error:', msg);
                vscode.window.showErrorMessage(`Failed to start server: ${msg}`);
            }
        },
        
        'aionmcp.stopServer': async () => {
            console.log('Stop command received');
            try {
                await serverManager.stopServer();
            } catch (error) {
                const msg = error instanceof Error ? error.message : String(error);
                console.error('Stop error:', msg);
                vscode.window.showErrorMessage(`Failed to stop server: ${msg}`);
            }
        },
        
        'aionmcp.restartServer': async () => {
            console.log('Restart command received');
            try {
                await serverManager.restartServer();
            } catch (error) {
                const msg = error instanceof Error ? error.message : String(error);
                console.error('Restart error:', msg);
                vscode.window.showErrorMessage(`Failed to restart server: ${msg}`);
            }
        },
        
        'aionmcp.refreshTools': () => {
            console.log('Refresh tools command');
            toolTreeProvider.refresh();
        },
        
        'aionmcp.refreshAgents': () => {
            console.log('Refresh agents command');
            agentTreeProvider.refresh();
        },
        
        'aionmcp.refreshAll': () => {
            console.log('Refresh all command');
            toolTreeProvider.refresh();
            agentTreeProvider.refresh();
            serverStatusProvider.refresh();
        },
        
        'aionmcp.executeTool': (item: ToolItem) => {
            console.log('Execute tool:', item?.tool?.name);
            if (item?.tool) {
                vscode.window.showInformationMessage(`Executing tool: ${item.tool.name}`);
            } else {
                vscode.window.showWarningMessage('No tool selected');
            }
        },
        
        'aionmcp.viewLogs': () => {
            console.log('View logs command');
            serverManager.showOutputChannel();
        },
        
        'aionmcp.openDashboard': () => {
            console.log('Open dashboard command');
            DashboardWebviewProvider.createOrShow(context.extensionUri, serverManager);
        },
        
        'aionmcp.openToolExecutor': (item: ToolItem) => {
            console.log('Open tool executor for:', item?.tool?.name);
            if (item?.tool) {
                ToolExecutorWebviewProvider.createOrShow(context.extensionUri, serverManager, item.tool);
            }
        }
    };

    for (const [command, handler] of Object.entries(commands)) {
        context.subscriptions.push(vscode.commands.registerCommand(command, handler));
    }
}