import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
    console.log('Test extension is now active!');
    
    // Register a simple test command
    let disposable = vscode.commands.registerCommand('aionmcp.testCommand', () => {
        vscode.window.showInformationMessage('Test command works!');
    });
    
    context.subscriptions.push(disposable);
    
    // Also register the start server command for testing
    let startServerDisposable = vscode.commands.registerCommand('aionmcp.startServer', () => {
        vscode.window.showInformationMessage('Start server command triggered!');
    });
    
    context.subscriptions.push(startServerDisposable);
    
    vscode.window.showInformationMessage('AionMCP test extension activated!');
}

export function deactivate() {}