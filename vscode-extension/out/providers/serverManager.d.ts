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
export declare class ServerManager implements vscode.Disposable {
    private context;
    private serverProcess;
    private isRunning;
    private readonly outputChannel;
    private readonly stateChangeEmitter;
    private readonly axiosInstance;
    readonly onServerStateChanged: vscode.Event<boolean>;
    constructor(context: vscode.ExtensionContext);
    private getWorkspaceRoot;
    private getConfig;
    private isGoProject;
    startServer(): Promise<void>;
    private handleServerError;
    private waitForServerReady;
    stopServer(): Promise<void>;
    private stopServerProcess;
    restartServer(): Promise<void>;
    private getApiUrl;
    isServerRunning(): boolean;
    getServerStatus(): {
        isRunning: boolean;
        port: number;
        grpcPort: number;
    };
    getDashboardData(): Promise<ServerStats>;
    getTools(): Promise<Tool[]>;
    getTool(toolName: string): Promise<Tool | undefined>;
    getAgents(): Promise<Agent[]>;
    executeTool(toolName: string, parameters: any): Promise<any>;
    getServerStats(): Promise<ServerStats | null>;
    private handleApiError;
    showOutputChannel(): void;
    importSpec(specType: string, path: string, name?: string): Promise<any>;
    registerMockAgent(): Promise<any>;
    dispose(): void;
}
//# sourceMappingURL=serverManager.d.ts.map