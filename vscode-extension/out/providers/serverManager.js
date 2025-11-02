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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ServerManager = void 0;
const vscode = __importStar(require("vscode"));
const cp = __importStar(require("child_process"));
const path = __importStar(require("path"));
const axios_1 = __importDefault(require("axios"));
class ServerManager {
    context;
    serverProcess = null;
    isRunning = false;
    config;
    outputChannel;
    stateChangeEmitter = new vscode.EventEmitter();
    onServerStateChanged = this.stateChangeEmitter.event;
    constructor(context) {
        this.context = context;
        this.config = vscode.workspace.getConfiguration('aionmcp');
        this.outputChannel = vscode.window.createOutputChannel('AionMCP Server');
        this.context.subscriptions.push(this.outputChannel);
    }
    async startServer() {
        if (this.isRunning) {
            throw new Error('Server is already running');
        }
        const serverPath = this.getServerPath();
        const workspaceRoot = this.getWorkspaceRoot();
        if (!serverPath) {
            throw new Error('AionMCP server path not configured');
        }
        this.outputChannel.appendLine(`Starting AionMCP server: ${serverPath}`);
        this.outputChannel.show();
        return new Promise((resolve, reject) => {
            const env = {
                ...process.env,
                AIONMCP_LOG_LEVEL: this.config.get('logLevel', 'info'),
                AIONMCP_HTTP_PORT: this.config.get('serverPort', 8080).toString(),
                AIONMCP_GRPC_PORT: this.config.get('grpcPort', 50051).toString()
            };
            this.serverProcess = cp.spawn(serverPath, [], {
                cwd: workspaceRoot,
                env,
                stdio: ['pipe', 'pipe', 'pipe']
            });
            this.serverProcess.stdout?.on('data', (data) => {
                this.outputChannel.appendLine(data.toString());
            });
            this.serverProcess.stderr?.on('data', (data) => {
                this.outputChannel.appendLine(`ERROR: ${data.toString()}`);
            });
            this.serverProcess.on('error', (error) => {
                this.outputChannel.appendLine(`Failed to start server: ${error.message}`);
                this.isRunning = false;
                this.stateChangeEmitter.fire(false);
                reject(error);
            });
            this.serverProcess.on('exit', (code) => {
                this.outputChannel.appendLine(`Server exited with code: ${code}`);
                this.isRunning = false;
                this.serverProcess = null;
                this.stateChangeEmitter.fire(false);
            });
            // Wait for server to be ready
            this.waitForServer().then(() => {
                this.isRunning = true;
                this.stateChangeEmitter.fire(true);
                this.outputChannel.appendLine('Server is ready and accepting connections');
                resolve();
            }).catch(reject);
        });
    }
    async stopServer() {
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
            }
            else {
                resolve();
            }
        });
    }
    async restartServer() {
        if (this.isRunning) {
            await this.stopServer();
            // Wait a bit before restarting
            await new Promise(resolve => setTimeout(resolve, 1000));
        }
        await this.startServer();
    }
    async getTools() {
        if (!this.isRunning) {
            return [];
        }
        try {
            const response = await axios_1.default.get(`http://localhost:${this.getServerPort()}/api/tools`);
            return response.data.tools || [];
        }
        catch (error) {
            console.error('Failed to fetch tools:', error);
            return [];
        }
    }
    async getTool(name) {
        if (!this.isRunning) {
            return null;
        }
        try {
            const response = await axios_1.default.get(`http://localhost:${this.getServerPort()}/api/tools/${encodeURIComponent(name)}`);
            return response.data;
        }
        catch (error) {
            console.error(`Failed to fetch tool ${name}:`, error);
            return null;
        }
    }
    async executeTool(toolName, args, context) {
        if (!this.isRunning) {
            throw new Error('Server is not running');
        }
        try {
            const response = await axios_1.default.post(`http://localhost:${this.getServerPort()}/api/tools/${encodeURIComponent(toolName)}/invoke`, {
                args,
                context: context || {}
            });
            return response.data;
        }
        catch (error) {
            throw new Error(`Tool execution failed: ${error.response?.data?.message || error.message}`);
        }
    }
    async getAgents() {
        if (!this.isRunning) {
            return [];
        }
        try {
            const response = await axios_1.default.get(`http://localhost:${this.getServerPort()}/api/agents`);
            return response.data.agents || [];
        }
        catch (error) {
            console.error('Failed to fetch agents:', error);
            return [];
        }
    }
    async getServerStats() {
        if (!this.isRunning) {
            return null;
        }
        try {
            const response = await axios_1.default.get(`http://localhost:${this.getServerPort()}/api/admin/stats`);
            return response.data;
        }
        catch (error) {
            console.error('Failed to fetch server stats:', error);
            return null;
        }
    }
    async importApiSpec(filePath) {
        if (!this.isRunning) {
            throw new Error('Server is not running');
        }
        try {
            // For now, just copy the file to the specs directory
            // In a real implementation, you'd POST the spec to the server
            const specDirs = this.config.get('specDirectories', ['./examples/specs']);
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
        }
        catch (error) {
            throw new Error(`Failed to import API spec: ${error.message}`);
        }
    }
    async refreshSpecs() {
        if (!this.isRunning) {
            return;
        }
        try {
            await axios_1.default.post(`http://localhost:${this.getServerPort()}/api/admin/reload-specs`);
        }
        catch (error) {
            // Ignore if endpoint doesn't exist yet
            console.log('Spec reload endpoint not available');
        }
    }
    getServerStatus() {
        return {
            isRunning: this.isRunning,
            port: this.getServerPort(),
            grpcPort: this.getGrpcPort()
        };
    }
    async waitForServer(maxRetries = 30, delayMs = 1000) {
        for (let i = 0; i < maxRetries; i++) {
            try {
                await axios_1.default.get(`http://localhost:${this.getServerPort()}/api/health`, {
                    timeout: 2000
                });
                return; // Server is ready
            }
            catch (error) {
                if (i === maxRetries - 1) {
                    throw new Error('Server failed to start within timeout period');
                }
                await new Promise(resolve => setTimeout(resolve, delayMs));
            }
        }
    }
    getServerPath() {
        const configured = this.config.get('serverPath');
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
    getWorkspaceRoot() {
        const workspaceFolders = vscode.workspace.workspaceFolders;
        if (workspaceFolders && workspaceFolders.length > 0) {
            return workspaceFolders[0].uri.fsPath;
        }
        return process.cwd();
    }
    getServerPort() {
        return this.config.get('serverPort', 8080);
    }
    getGrpcPort() {
        return this.config.get('grpcPort', 50051);
    }
    dispose() {
        if (this.isRunning && this.serverProcess) {
            this.serverProcess.kill();
        }
        this.stateChangeEmitter.dispose();
    }
}
exports.ServerManager = ServerManager;
//# sourceMappingURL=serverManager.js.map