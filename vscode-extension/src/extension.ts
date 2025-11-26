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
        },

        'aionmcp.importApiSpec': async () => {
            console.log('Import spec command received');
            try {
                // Ask user for spec type
                const specType = await vscode.window.showQuickPick(
                    ['openapi', 'graphql', 'asyncapi'],
                    { placeHolder: 'Select specification type' }
                );
                if (!specType) {
                    return;
                }

                // Ask for file path or URL
                const path = await vscode.window.showInputBox({
                    prompt: 'Enter file path or URL to the specification',
                    placeHolder: 'e.g., ./specs/api.yaml or https://api.example.com/openapi.json'
                });
                if (!path) {
                    return;
                }

                // Ask for optional name
                const name = await vscode.window.showInputBox({
                    prompt: 'Optional: Enter a name for this specification',
                    placeHolder: 'e.g., My API'
                });

                // Import the spec
                const result = await serverManager.importSpec(specType, path, name || '');
                
                vscode.window.showInformationMessage(
                    `Successfully imported ${result.tools?.length || 0} tools from ${specType} spec`
                );
                
                // Refresh tools tree
                toolTreeProvider.refresh();
            } catch (error) {
                const msg = error instanceof Error ? error.message : String(error);
                console.error('Import spec error:', msg);
                vscode.window.showErrorMessage(`Failed to import specification: ${msg}`);
            }
        },

        'aionmcp.registerCopilotAgent': async () => {
            console.log('Register Copilot agent command received');
            try {
                if (!serverManager.isServerRunning()) {
                    vscode.window.showErrorMessage('Server is not running. Please start the server first.');
                    return;
                }

                await serverManager.registerMockAgent();
                vscode.window.showInformationMessage('GitHub Copilot agent registered successfully!');
                
                // Refresh agents tree
                agentTreeProvider.refresh();
            } catch (error) {
                const msg = error instanceof Error ? error.message : String(error);
                console.error('Register agent error:', msg);
                vscode.window.showErrorMessage(`Failed to register Copilot agent: ${msg}`);
            }
        },

        'aionmcp.showAPIInfo': async () => {
            console.log('Show API info command received');
            try {
                const baseUrl = `http://localhost:8080`;
                
                const message = 
                    `AionMCP Server is running at ${baseUrl}\n\n` +
                    `API Endpoints:\n` +
                    `• GET  ${baseUrl}/api/v1/tools - List all tools\n` +
                    `• GET  ${baseUrl}/api/v1/tools/{id} - Get tool details\n` +
                    `• POST ${baseUrl}/api/v1/tools/{id}/invoke - Execute tool\n\n` +
                    `Use these endpoints from any agent or application!\n` +
                    `See documentation for integration examples.`;
                
                const action = await vscode.window.showInformationMessage(
                    message,
                    'Copy API URL',
                    'View Documentation',
                    'Done'
                );

                if (action === 'Copy API URL') {
                    await vscode.env.clipboard.writeText(baseUrl);
                    vscode.window.showInformationMessage(`Copied ${baseUrl} to clipboard`);
                } else if (action === 'View Documentation') {
                    // Open documentation in browser
                    vscode.env.openExternal(vscode.Uri.parse('https://github.com/kiransth77/aionmcp#model-independent-tool-server'));
                }
            } catch (error) {
                const msg = error instanceof Error ? error.message : String(error);
                console.error('Show API info error:', msg);
                vscode.window.showErrorMessage(`Failed to show API info: ${msg}`);
            }
        }
    };

    for (const [command, handler] of Object.entries(commands)) {
        context.subscriptions.push(vscode.commands.registerCommand(command, handler));
    }
}