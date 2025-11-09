import * as vscode from 'vscode';
export interface Tool {
    name: string;
    description: string;
    inputSchema: {
        type: string;
        properties: Record<string, any>;
        required?: string[];
    };
    version?: string;
    source?: string;
}
export interface Agent {
    id: string;
    name: string;
    capabilities: string[];
    status: 'connected' | 'disconnected';
    lastSeen: Date;
}
export interface ServerStats {
    uptime: number;
    toolCount: number;
    agentCount: number;
    executionCount: number;
    successRate: number;
}
export declare class ServerManager {
    private context;
    private serverProcess;
    private isRunning;
    private config;
    private outputChannel;
    private stateChangeEmitter;
    readonly onServerStateChanged: vscode.Event<boolean>;
    constructor(context: vscode.ExtensionContext);
    startServer(): Promise<void>;
    stopServer(): Promise<void>;
    restartServer(): Promise<void>;
    getTools(): Promise<Tool[]>;
    getTool(name: string): Promise<Tool | null>;
    executeTool(toolName: string, args: any, context?: any): Promise<any>;
    getAgents(): Promise<Agent[]>;
    getServerStats(): Promise<ServerStats | null>;
    importApiSpec(filePath: string): Promise<void>;
    refreshSpecs(): Promise<void>;
    getServerStatus(): {
        isRunning: boolean;
        port: number;
        grpcPort: number;
    };
    private waitForServer;
    private getServerPath;
    private getWorkspaceRoot;
    private getServerPort;
    private getGrpcPort;
    dispose(): void;
}
//# sourceMappingURL=serverManager.d.ts.map