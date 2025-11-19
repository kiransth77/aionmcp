import * as vscode from 'vscode';
import * as cp from 'child_process';
import * as path from 'path';
import axios from 'axios';
import { promisify } from 'util';
import * as fs from 'fs';

const exec = promisify(cp.exec);

// Interfaces remain the same
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

export class ServerManager implements vscode.Disposable {
    private serverProcess: cp.ChildProcess | null = null;
    private isRunning = false;
    private readonly outputChannel: vscode.OutputChannel;
    private readonly stateChangeEmitter = new vscode.EventEmitter<boolean>();
    
    public readonly onServerStateChanged = this.stateChangeEmitter.event;
    
    constructor(private context: vscode.ExtensionContext) {
        this.outputChannel = vscode.window.createOutputChannel('AionMCP Server');
        this.context.subscriptions.push(this.outputChannel, this);
    }

    private getWorkspaceRoot(): string | undefined {
        const workspaceFolders = vscode.workspace.workspaceFolders;
        console.log('Workspace folders:', workspaceFolders);
        
        if (!workspaceFolders || workspaceFolders.length === 0) {
            console.warn('No workspace folders are open');
            return undefined;
        }
        
        const root = workspaceFolders[0].uri.fsPath;
        console.log('Using workspace root:', root);
        return root;
    }

    private getConfig<T>(key: string, defaultValue: T): T {
        return vscode.workspace.getConfiguration('aionmcp').get<T>(key, defaultValue);
    }

    private async isGoProject(workspaceRoot: string): Promise<boolean> {
        try {
            await fs.promises.stat(path.join(workspaceRoot, 'go.mod'));
            return true;
        } catch {
            return false;
        }
    }

    async startServer(): Promise<void> {
        console.log('[ServerManager] startServer called');
        
        if (this.isRunning) {
            console.log('[ServerManager] Server already running');
            vscode.window.showInformationMessage('AionMCP server is already running.');
            return;
        }

        console.log('[ServerManager] Getting workspace root...');
        const workspaceRoot = this.getWorkspaceRoot();
        
        if (!workspaceRoot) {
            const message = 'No workspace folder is open. Please open the aionmcp project folder in VS Code.';
            console.error('[ServerManager]', message);
            throw new Error(message);
        }

        console.log('[ServerManager] Checking if Go project...');
        if (!await this.isGoProject(workspaceRoot)) {
            const message = 'The current workspace does not appear to be a Go project (go.mod not found).';
            console.error('[ServerManager]', message);
            throw new Error(message);
        }

        console.log('[ServerManager] Clearing output channel...');
        this.outputChannel.clear();
        this.outputChannel.show();
        this.outputChannel.appendLine('Starting AionMCP server...');

        const serverPort = this.getConfig('serverPort', 8080);
        const grpcPort = this.getConfig('grpcPort', 50051);
        const logLevel = this.getConfig('logLevel', 'info');

        const env = {
            ...process.env,
            AIONMCP_HTTP_PORT: String(serverPort),
            AIONMCP_GRPC_PORT: String(grpcPort),
            AIONMCP_LOG_LEVEL: logLevel,
        };

        this.outputChannel.appendLine(`Workspace: ${workspaceRoot}`);
        this.outputChannel.appendLine(`Command: go run cmd/server/main.go`);
        this.outputChannel.appendLine(`Environment: HTTP_PORT=${env.AIONMCP_HTTP_PORT}, GRPC_PORT=${env.AIONMCP_GRPC_PORT}, LOG_LEVEL=${env.AIONMCP_LOG_LEVEL}`);

        console.log('[ServerManager] Spawning go process...');
        this.serverProcess = cp.spawn('go', ['run', 'cmd/server/main.go'], {
            cwd: workspaceRoot,
            env,
        });

        console.log('[ServerManager] Setting up process listeners...');
        this.serverProcess.stdout?.on('data', (data) => {
            const output = data.toString();
            console.log('[ServerProcess stdout]', output);
            this.outputChannel.append(output);
        });
        
        this.serverProcess.stderr?.on('data', (data) => {
            const output = data.toString();
            console.log('[ServerProcess stderr]', output);
            this.outputChannel.append(output);
        });

        this.serverProcess.on('error', (err) => {
            console.error('[ServerProcess error event]', err);
            this.handleServerError(new Error(`Failed to start server process: ${err.message}`));
        });

        this.serverProcess.on('exit', (code) => {
            console.log('[ServerProcess exit]', code);
            if (this.isRunning) { // If it was running, this is an unexpected exit
                this.handleServerError(new Error(`Server process exited unexpectedly with code ${code}.`));
            } else { // Normal exit
                this.outputChannel.appendLine(`Server process stopped with code ${code}.`);
            }
        });

        try {
            console.log('[ServerManager] Waiting for server to be ready...');
            await this.waitForServerReady(serverPort);
            this.isRunning = true;
            this.stateChangeEmitter.fire(true);
            this.outputChannel.appendLine('[SUCCESS] AionMCP server is running and connected.');
            vscode.window.showInformationMessage('AionMCP server started successfully.');
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : String(error);
            console.error('[ServerManager] Server startup failed:', errorMessage);
            this.handleServerError(new Error(`Server failed to become ready: ${errorMessage}`));
            this.stopServerProcess();
            throw error; // Re-throw to notify the caller
        }
    }

    private handleServerError(error: Error) {
        this.outputChannel.appendLine(`[ERROR] ${error.message}`);
        vscode.window.showErrorMessage(`AionMCP Server Error: ${error.message}`);
        this.isRunning = false;
        this.stateChangeEmitter.fire(false);
    }

    private async waitForServerReady(port: number, timeout = 15000): Promise<void> {
        this.outputChannel.appendLine(`[HEALTH] Waiting for server to be ready on port ${port}...`);
        const startTime = Date.now();
        
        while (Date.now() - startTime < timeout) {
            try {
                // The server might not have a dedicated health endpoint yet, so we check the base URL.
                await axios.get(`http://localhost:${port}/`, { timeout: 1000 });
                this.outputChannel.appendLine('[HEALTH] Server is responsive.');
                return;
            } catch (error) {
                await new Promise(resolve => setTimeout(resolve, 1000)); // Wait before retrying
            }
        }
        throw new Error(`Server did not respond on port ${port} within ${timeout / 1000} seconds.`);
    }

    async stopServer(): Promise<void> {
        if (!this.isRunning) {
            vscode.window.showInformationMessage('AionMCP server is not running.');
            return;
        }
        this.outputChannel.appendLine('Stopping AionMCP server...');
        await this.stopServerProcess();
        vscode.window.showInformationMessage('AionMCP server stopped.');
    }

    private stopServerProcess(): Promise<void> {
        return new Promise((resolve) => {
            if (!this.serverProcess) {
                this.isRunning = false;
                this.stateChangeEmitter.fire(false);
                return resolve();
            }

            this.serverProcess.once('exit', () => {
                this.isRunning = false;
                this.serverProcess = null;
                this.stateChangeEmitter.fire(false);
                this.outputChannel.appendLine('Server process terminated.');
                resolve();
            });

            this.serverProcess.kill('SIGTERM');

            // Failsafe
            setTimeout(() => {
                if (this.serverProcess) {
                    this.outputChannel.appendLine('Server did not respond to SIGTERM, forcing shutdown with SIGKILL.');
                    this.serverProcess.kill('SIGKILL');
                }
            }, 3000);
        });
    }

    async restartServer(): Promise<void> {
        this.outputChannel.appendLine('Restarting AionMCP server...');
        if (this.isRunning) {
            await this.stopServerProcess();
        }
        await this.startServer();
    }

    private getApiUrl(): string {
        const port = this.getConfig('serverPort', 8080);
        return `http://localhost:${port}/api/v1`;
    }

    public isServerRunning(): boolean {
        return this.isRunning;
    }

    public getServerStatus(): { isRunning: boolean; port: number; grpcPort: number } {
        return {
            isRunning: this.isRunning,
            port: this.getConfig('serverPort', 8080),
            grpcPort: this.getConfig('grpcPort', 50051)
        };
    }

    async getDashboardData(): Promise<ServerStats> {
        if (!this.isRunning) {
            throw new Error('Server is not running.');
        }
        const port = this.getConfig('serverPort', 8080);
        try {
            const response = await axios.get(`http://localhost:${port}/api/v1/stats`);
            return response.data;
        } catch (error) {
            if (axios.isAxiosError(error)) {
                throw new Error(`Failed to fetch dashboard data from server: ${error.message}`);
            }
            throw error;
        }
    }

    async getTools(): Promise<Tool[]> {
        if (!this.isRunning) {
            return [];
        }
        const port = this.getConfig('serverPort', 8080);
        try {
            const response = await axios.get(`http://localhost:${port}/api/v1/tools`);
            return response.data.tools || [];
        } catch (error) {
            console.error('Failed to get tools:', error);
            return [];
        }
    }

    async getTool(toolName: string): Promise<Tool | undefined> {
        if (!this.isRunning) {
            return undefined;
        }
        const port = this.getConfig('serverPort', 8080);
        try {
            const response = await axios.get(`http://localhost:${port}/api/v1/tools/${toolName}`);
            return response.data;
        } catch (error) {
            console.error(`Failed to get tool ${toolName}:`, error);
            return undefined;
        }
    }

    async getAgents(): Promise<Agent[]> {
        if (!this.isRunning) {
            return [];
        }
        const port = this.getConfig('serverPort', 8080);
        try {
            const response = await axios.get(`http://localhost:${port}/api/v1/agents`);
            return response.data.agents || [];
        } catch (error) {
            console.error('Failed to get agents:', error);
            return [];
        }
    }

    async executeTool(toolName: string, parameters: any): Promise<any> {
        if (!this.isRunning) throw new Error('Server is not running.');
        try {
            const response = await axios.post(`${this.getApiUrl()}/tools/${toolName}/execute`, { args: parameters });
            return response.data;
        } catch (error) {
            this.handleApiError(`Execution of tool '${toolName}' failed`, error);
            throw error;
        }
    }
    
    async getServerStats(): Promise<ServerStats | null> {
        if (!this.isRunning) return null;
        try {
            const response = await axios.get(`${this.getApiUrl()}/stats`);
            return response.data;
        } catch (error) {
            this.handleApiError('Failed to fetch server stats', error);
            return null;
        }
    }

    private handleApiError(message: string, error: any) {
        const errorMessage = axios.isAxiosError(error) && error.response?.data?.error
            ? error.response.data.error
            : error.message;
        vscode.window.showErrorMessage(`${message}: ${errorMessage}`);
        this.outputChannel.appendLine(`[API ERROR] ${message}: ${errorMessage}`);
    }

    public showOutputChannel(): void {
        this.outputChannel.show();
    }

    dispose(): void {
        this.stopServerProcess();
        this.outputChannel.dispose();
        this.stateChangeEmitter.dispose();
    }
}