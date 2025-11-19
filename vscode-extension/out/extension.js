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
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = __importStar(require("vscode"));
const serverManager_1 = require("./providers/serverManager");
const toolTreeProvider_1 = require("./providers/toolTreeProvider");
const agentTreeProvider_1 = require("./providers/agentTreeProvider");
const serverStatusProvider_1 = require("./providers/serverStatusProvider");
const toolExecutorWebview_1 = require("./webviews/toolExecutorWebview");
const logOutputProvider_1 = require("./providers/logOutputProvider");
const dashboardWebview_1 = require("./webviews/dashboardWebview");
let serverManager;
let toolTreeProvider;
let agentTreeProvider;
let serverStatusProvider;
let logOutputProvider;
function activate(context) {
    console.log('AionMCP extension is now active!');
    vscode.window.showInformationMessage('AionMCP extension activated successfully!');
    // Initialize providers
    serverManager = new serverManager_1.ServerManager(context);
    toolTreeProvider = new toolTreeProvider_1.ToolTreeProvider(serverManager);
    agentTreeProvider = new agentTreeProvider_1.AgentTreeProvider(serverManager);
    serverStatusProvider = new serverStatusProvider_1.ServerStatusProvider(serverManager);
    logOutputProvider = new logOutputProvider_1.LogOutputProvider();
    // Register tree data providers
    vscode.window.registerTreeDataProvider('aionmcp.toolsView', toolTreeProvider);
    vscode.window.registerTreeDataProvider('aionmcp.agentsView', agentTreeProvider);
    vscode.window.registerTreeDataProvider('aionmcp.serverView', serverStatusProvider);
    // Register commands
    registerCommands(context);
    // Set extension as active
    vscode.commands.executeCommand('setContext', 'aionmcp.extensionActive', true);
    // Auto-start server if configured
    const config = vscode.workspace.getConfiguration('aionmcp');
    if (config.get('autoStart', false)) {
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
    serverManager.onServerStateChanged((isRunning) => {
        vscode.commands.executeCommand('setContext', 'aionmcp.serverRunning', isRunning);
        statusBarItem.text = isRunning
            ? '$(server-process) AionMCP $(check)'
            : '$(server-process) AionMCP $(x)';
        statusBarItem.color = isRunning ? undefined : new vscode.ThemeColor('statusBarItem.errorForeground');
    });
}
// In extension.ts or a dedicated command registration file
function registerCommands(context) {
    const commands = {
        'aionmcp.startServer': async () => {
            try {
                await serverManager.startServer();
            }
            catch (error) {
                const message = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(`Failed to start AionMCP server: ${message}`);
            }
        },
        'aionmcp.stopServer': async () => {
            try {
                await serverManager.stopServer();
            }
            catch (error) {
                const message = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(`Failed to stop AionMCP server: ${message}`);
            }
        },
        'aionmcp.restartServer': async () => {
            try {
                await serverManager.restartServer();
            }
            catch (error) {
                const message = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(`Failed to restart AionMCP server: ${message}`);
            }
        },
        'aionmcp.refreshAll': () => {
            toolTreeProvider.refresh();
            agentTreeProvider.refresh();
            serverStatusProvider.refresh();
        },
        'aionmcp.executeTool': (item) => {
            if (item?.tool) {
                vscode.commands.executeCommand('aionmcp.openToolExecutor', item);
            }
            else {
                vscode.window.showWarningMessage('No tool selected to execute.');
            }
        },
        'aionmcp.openToolExecutor': (item) => {
            toolExecutorWebview_1.ToolExecutorWebviewProvider.createOrShow(context.extensionUri, serverManager, item?.tool);
        },
        'aionmcp.showDashboard': () => {
            dashboardWebview_1.DashboardWebviewProvider.createOrShow(context.extensionUri, serverManager);
        },
        'aionmcp.viewLogs': () => {
            // This command is now handled by focusing the output channel, see below
            vscode.commands.executeCommand('aionmcp.output.focus');
        }
    };
    for (const [command, handler] of Object.entries(commands)) {
        context.subscriptions.push(vscode.commands.registerCommand(command, handler));
    }
    // Special handling for the output channel focus
    context.subscriptions.push(vscode.commands.registerCommand('aionmcp.output.focus', () => {
        // The ServerManager now creates and manages the output channel.
        // We can ask it to show the channel.
        // This requires a new method in ServerManager: `showOutputChannel()`
        serverManager.showOutputChannel();
    }));
}
async function executeToolQuick(toolItem) {
    // Simple tool execution with basic parameter input
    const tool = toolItem.tool;
    // If tool has no parameters, execute directly
    if (!tool.inputSchema || Object.keys(tool.inputSchema.properties || {}).length === 0) {
        const result = await serverManager.executeTool(tool.name, {});
        showExecutionResult(tool.name, result);
        return true;
    }
    // For tools with parameters, show quick input
    const params = {};
    const properties = tool.inputSchema.properties || {};
    for (const [propName, propSchema] of Object.entries(properties)) {
        const schema = propSchema;
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
        }
        else if (schema.type === 'boolean') {
            params[propName] = value.toLowerCase() === 'true';
        }
        else {
            params[propName] = value;
        }
    }
    const result = await serverManager.executeTool(tool.name, params);
    showExecutionResult(tool.name, result);
    return true;
}
function showExecutionResult(toolName, result) {
    const resultStr = typeof result === 'string' ? result : JSON.stringify(result, null, 2);
    // Show result in output channel
    logOutputProvider.appendLine(`\n=== Tool Execution: ${toolName} ===`);
    logOutputProvider.appendLine(resultStr);
    logOutputProvider.appendLine('=== End Result ===\n');
    logOutputProvider.show();
}
function deactivate() {
    if (serverManager) {
        serverManager.dispose();
    }
}
//# sourceMappingURL=extension.js.map