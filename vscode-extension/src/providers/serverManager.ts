import * as vscode from 'vscode';
import * as cp from 'child_process';
import * as path from 'path';
import axios from 'axios';

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

export class ServerManager {
    private serverProcess: cp.ChildProcess | null = null;
    private isRunning = false;
    private config: vscode.WorkspaceConfiguration;
    private outputChannel: vscode.OutputChannel;
    private stateChangeEmitter = new vscode.EventEmitter<boolean>();
    
    public readonly onServerStateChanged = this.stateChangeEmitter.event;
    
    constructor(private context: vscode.ExtensionContext) {
        this.config = vscode.workspace.getConfiguration('aionmcp');
        this.outputChannel = vscode.window.createOutputChannel('AionMCP Server');
        this.context.subscriptions.push(this.outputChannel);
    }
    
    async startServer(): Promise<void> {
        if (this.isRunning) {
            throw new Error('Server is already running');
        }
        
        const serverPath = this.getServerPath();
        const workspaceRoot = this.getWorkspaceRoot();
        
        if (!serverPath) {
            throw new Error('AionMCP server path not configured');
        }
        
        // Check if binary exists
        const fs = require('fs');
        if (!fs.existsSync(serverPath)) {
            throw new Error(`AionMCP server binary not found at: ${serverPath}`);
        }
        
        this.outputChannel.appendLine(`Starting AionMCP server: ${serverPath}`);
        this.outputChannel.appendLine(`Working directory: ${workspaceRoot}`);
        this.outputChannel.show();
        
        return new Promise((resolve, reject) => {
            const env = {
                ...process.env,
                AIONMCP_LOG_LEVEL: this.config.get<string>('logLevel', 'info'),
                AIONMCP_HTTP_PORT: this.config.get<number>('serverPort', 8080).toString(),
                AIONMCP_GRPC_PORT: this.config.get<number>('grpcPort', 50051).toString()
            };
            
            this.outputChannel.appendLine(`Environment: HTTP_PORT=${env.AIONMCP_HTTP_PORT}, GRPC_PORT=${env.AIONMCP_GRPC_PORT}, LOG_LEVEL=${env.AIONMCP_LOG_LEVEL}`);
            
            this.serverProcess = cp.spawn(serverPath, [], {
                cwd: workspaceRoot,
                env,
                stdio: ['pipe', 'pipe', 'pipe']
            });
            
            this.serverProcess.stdout?.on('data', (data) => {
                this.outputChannel.appendLine(`[STDOUT] ${data.toString()}`);
            });
            
            this.serverProcess.stderr?.on('data', (data) => {
                this.outputChannel.appendLine(`[STDERR] ${data.toString()}`);
            });
            
            this.serverProcess.on('error', (error) => {
                this.outputChannel.appendLine(`[ERROR] Failed to start server: ${error.message}`);
                this.isRunning = false;
                this.stateChangeEmitter.fire(false);
                reject(new Error(`Failed to start server: ${error.message}`));
            });
            
            this.serverProcess.on('exit', (code, signal) => {
                this.outputChannel.appendLine(`[EXIT] Server exited with code: ${code}, signal: ${signal}`);
                this.isRunning = false;
                this.serverProcess = null;
                this.stateChangeEmitter.fire(false);
                
                if (code !== 0 && code !== null) {
                    reject(new Error(`Server exited with non-zero code: ${code}`));
                }
            });
            
            // Wait for server to be ready
            this.waitForServer().then(() => {
                this.isRunning = true;
                this.stateChangeEmitter.fire(true);
                this.outputChannel.appendLine('[SUCCESS] Server is ready and accepting connections');
                resolve();
            }).catch((error) => {
                this.outputChannel.appendLine(`[TIMEOUT] Server failed to start within timeout: ${error.message}`);
                reject(error);
            });
        });
    }
    
    async stopServer(): Promise<void> {
        if (!this.isRunning || !this.serverProcess) {
            throw new Error('Server is not running');
        }
        
        this.outputChannel.appendLine('Stopping AionMCP server...');
        
        return new Promise((resolve) => {
            if (this.serverProcess) {
                this.serverProcess.on('exit', () => {
                    this.outputChannel.appendLine('Server stopped');
                    resolve();
                });
                
                // Try graceful shutdown first
                this.serverProcess.kill('SIGTERM');
                
                // Force kill after 5 seconds
                setTimeout(() => {
                    if (this.serverProcess && !this.serverProcess.killed) {
                        this.serverProcess.kill('SIGKILL');
                    }
                }, 5000);
            } else {
                resolve();
            }
        });
    }
    
    async restartServer(): Promise<void> {
        if (this.isRunning) {
            await this.stopServer();
            // Wait a bit before restarting
            await new Promise(resolve => setTimeout(resolve, 1000));
        }
        await this.startServer();
    }
    
    async getTools(): Promise<Tool[]> {
        if (!this.isRunning) {
            return [];
        }
        
        try {
            const response = await axios.get(`http://localhost:${this.getServerPort()}/api/tools`);
            return response.data.tools || [];
        } catch (error) {
            console.error('Failed to fetch tools:', error);
            return [];
        }
    }
    
    async getTool(name: string): Promise<Tool | null> {
        if (!this.isRunning) {
            return null;
        }
        
        try {
            const response = await axios.get(`http://localhost:${this.getServerPort()}/api/tools/${encodeURIComponent(name)}`);
            return response.data;
        } catch (error) {
            console.error(`Failed to fetch tool ${name}:`, error);
            return null;
        }
    }
    
    async executeTool(toolName: string, args: any, context?: any): Promise<any> {
        if (!this.isRunning) {
            throw new Error('Server is not running');
        }
        
        try {
            const response = await axios.post(
                `http://localhost:${this.getServerPort()}/api/tools/${encodeURIComponent(toolName)}/invoke`,
                {
                    args,
                    context: context || {}
                }
            );
            return response.data;
        } catch (error: any) {
            throw new Error(`Tool execution failed: ${error.response?.data?.message || error.message}`);
        }
    }
    
    async getAgents(): Promise<Agent[]> {
        if (!this.isRunning) {
            return [];
        }
        
        try {
            const response = await axios.get(`http://localhost:${this.getServerPort()}/api/agents`);
            return response.data.agents || [];
        } catch (error) {
            console.error('Failed to fetch agents:', error);
            return [];
        }
    }
    
    async getServerStats(): Promise<ServerStats | null> {
        if (!this.isRunning) {
            return null;
        }
        
        try {
            const response = await axios.get(`http://localhost:${this.getServerPort()}/api/admin/stats`);
            return response.data;
        } catch (error) {
            console.error('Failed to fetch server stats:', error);
            return null;
        }
    }
    
    async importApiSpec(filePath: string): Promise<void> {
        if (!this.isRunning) {
            throw new Error('Server is not running');
        }
        
        try {
            // For now, just copy the file to the specs directory
            // In a real implementation, you'd POST the spec to the server
            const specDirs = this.config.get<string[]>('specDirectories', ['./examples/specs']);
            const targetDir = path.resolve(this.getWorkspaceRoot(), specDirs[0]);
            const fileName = path.basename(filePath);
            const targetPath = path.join(targetDir, fileName);
            
            const fs = require('fs');
            if (!fs.existsSync(targetDir)) {
                fs.mkdirSync(targetDir, { recursive: true });
            }
            
            fs.copyFileSync(filePath, targetPath);
            this.outputChannel.appendLine(`Imported API spec: ${fileName}`);
            
            // Trigger server reload if it supports hot reload
            await this.refreshSpecs();
        } catch (error: any) {
            throw new Error(`Failed to import API spec: ${error.message}`);
        }
    }
    
    async refreshSpecs(): Promise<void> {
        if (!this.isRunning) {
            return;
        }
        
        try {
            await axios.post(`http://localhost:${this.getServerPort()}/api/admin/reload-specs`);
        } catch (error) {
            // Ignore if endpoint doesn't exist yet
            console.log('Spec reload endpoint not available');
        }
    }
    
    getServerStatus(): { isRunning: boolean; port: number; grpcPort: number } {
        return {
            isRunning: this.isRunning,
            port: this.getServerPort(),
            grpcPort: this.getGrpcPort()
        };
    }
    
    private async waitForServer(maxRetries = 30, delayMs = 1000): Promise<void> {
        this.outputChannel.appendLine(`[HEALTH] Waiting for server to be ready (${maxRetries} retries, ${delayMs}ms delay)...`);
        
        for (let i = 0; i < maxRetries; i++) {
            try {
                const response = await axios.get(`http://localhost:${this.getServerPort()}/api/health`, {
                    timeout: 2000
                });
                this.outputChannel.appendLine(`[HEALTH] Server responded successfully: ${JSON.stringify(response.data)}`);
                return; // Server is ready
            } catch (error: any) {
                this.outputChannel.appendLine(`[HEALTH] Attempt ${i + 1}/${maxRetries} failed: ${error.message}`);
                if (i === maxRetries - 1) {
                    throw new Error(`Server failed to start within timeout period (${maxRetries * delayMs}ms). Check server logs for details.`);
                }
                await new Promise(resolve => setTimeout(resolve, delayMs));
            }
        }
    }
    
    private getServerPath(): string {
        const configured = this.config.get<string>('serverPath');
        if (configured && configured.trim()) {
            if (path.isAbsolute(configured)) {
                return configured;
            }
            return path.resolve(this.getWorkspaceRoot(), configured);
        }
        
        // Use bundled binary from extension
        const platform = process.platform;
        const extension = platform === 'win32' ? '.exe' : '';
        const binaryName = `aionmcp${extension}`;
        
        // Try extension's bundled binary first
        const bundledPath = path.join(this.context.extensionPath, 'bin', binaryName);
        const fs = require('fs');
        
        if (fs.existsSync(bundledPath)) {
            return bundledPath;
        }
        
        // Fallback to workspace bin directory
        const workspaceRoot = this.getWorkspaceRoot();
        return path.join(workspaceRoot, 'bin', binaryName);
    }
    
    private getWorkspaceRoot(): string {
        const workspaceFolders = vscode.workspace.workspaceFolders;
        if (workspaceFolders && workspaceFolders.length > 0) {
            return workspaceFolders[0].uri.fsPath;
        }
        return process.cwd();
    }
    
    private getServerPort(): number {
        return this.config.get<number>('serverPort', 8080);
    }
    
    private getGrpcPort(): number {
        return this.config.get<number>('grpcPort', 50051);
    }
    
    dispose(): void {
        if (this.isRunning && this.serverProcess) {
            this.serverProcess.kill();
        }
        this.stateChangeEmitter.dispose();
    }
}