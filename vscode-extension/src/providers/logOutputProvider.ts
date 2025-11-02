import * as vscode from 'vscode';

export class LogOutputProvider {
    private outputChannel: vscode.OutputChannel;
    
    constructor() {
        this.outputChannel = vscode.window.createOutputChannel('AionMCP Extension');
    }
    
    appendLine(message: string): void {
        const timestamp = new Date().toLocaleTimeString();
        this.outputChannel.appendLine(`[${timestamp}] ${message}`);
    }
    
    append(message: string): void {
        this.outputChannel.append(message);
    }
    
    clear(): void {
        this.outputChannel.clear();
    }
    
    show(): void {
        this.outputChannel.show();
    }
    
    hide(): void {
        this.outputChannel.hide();
    }
    
    dispose(): void {
        this.outputChannel.dispose();
    }
}