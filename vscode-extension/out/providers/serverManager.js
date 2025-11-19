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
const util_1 = require("util");
const fs = __importStar(require("fs"));
const exec = (0, util_1.promisify)(cp.exec);
class ServerManager {
    context;
    serverProcess = null;
    isRunning = false;
    outputChannel;
    stateChangeEmitter = new vscode.EventEmitter();
    onServerStateChanged = this.stateChangeEmitter.event;
    constructor(context) {
        this.context = context;
        this.outputChannel = vscode.window.createOutputChannel('AionMCP Server');
        this.context.subscriptions.push(this.outputChannel, this);
    }
    getWorkspaceRoot() {
        return vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
    }
    getConfig(key, defaultValue) {
        return vscode.workspace.getConfiguration('aionmcp').get(key, defaultValue);
    }
    async isGoProject(workspaceRoot) {
        try {
            await fs.promises.stat(path.join(workspaceRoot, 'go.mod'));
            return true;
        }
        catch {
            return false;
        }
    }
    async startServer() {
        if (this.isRunning) {
            vscode.window.showInformationMessage('AionMCP server is already running.');
            return;
        }
        const workspaceRoot = this.getWorkspaceRoot();
        if (!workspaceRoot) {
            throw new Error('No workspace folder is open.');
        }
        if (!await this.isGoProject(workspaceRoot)) {
            throw new Error('The current workspace does not appear to be a Go project (go.mod not found).');
        }
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
        this.serverProcess = cp.spawn('go', ['run', 'cmd/server/main.go'], {
            cwd: workspaceRoot,
            env,
        });
        this.serverProcess.stdout?.on('data', (data) => this.outputChannel.append(data.toString()));
        this.serverProcess.stderr?.on('data', (data) => this.outputChannel.append(data.toString()));
        this.serverProcess.on('error', (err) => {
            this.handleServerError(new Error(`Failed to start server process: ${err.message}`));
        });
        this.serverProcess.on('exit', (code) => {
            if (this.isRunning) { // If it was running, this is an unexpected exit
                this.handleServerError(new Error(`Server process exited unexpectedly with code ${code}.`));
            }
            else { // Normal exit
                this.outputChannel.appendLine(`Server process stopped with code ${code}.`);
            }
        });
        try {
            await this.waitForServerReady(serverPort);
            this.isRunning = true;
            this.stateChangeEmitter.fire(true);
            this.outputChannel.appendLine('[SUCCESS] AionMCP server is running and connected.');
            vscode.window.showInformationMessage('AionMCP server started successfully.');
        }
        catch (error) {
            const errorMessage = error instanceof Error ? error.message : String(error);
            this.handleServerError(new Error(`Server failed to become ready: ${errorMessage}`));
            this.stopServerProcess();
            throw error; // Re-throw to notify the caller
        }
    }
    handleServerError(error) {
        this.outputChannel.appendLine(`[ERROR] ${error.message}`);
        vscode.window.showErrorMessage(`AionMCP Server Error: ${error.message}`);
        this.isRunning = false;
        this.stateChangeEmitter.fire(false);
    }
    async waitForServerReady(port, timeout = 15000) {
        this.outputChannel.appendLine(`[HEALTH] Waiting for server to be ready on port ${port}...`);
        const startTime = Date.now();
        while (Date.now() - startTime < timeout) {
            try {
                // The server might not have a dedicated health endpoint yet, so we check the base URL.
                await axios_1.default.get(`http://localhost:${port}/`, { timeout: 1000 });
                this.outputChannel.appendLine('[HEALTH] Server is responsive.');
                return;
            }
            catch (error) {
                await new Promise(resolve => setTimeout(resolve, 1000)); // Wait before retrying
            }
        }
        throw new Error(`Server did not respond on port ${port} within ${timeout / 1000} seconds.`);
    }
    async stopServer() {
        if (!this.isRunning) {
            vscode.window.showInformationMessage('AionMCP server is not running.');
            return;
        }
        this.outputChannel.appendLine('Stopping AionMCP server...');
        await this.stopServerProcess();
        vscode.window.showInformationMessage('AionMCP server stopped.');
    }
    stopServerProcess() {
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
    async restartServer() {
        this.outputChannel.appendLine('Restarting AionMCP server...');
        if (this.isRunning) {
            await this.stopServerProcess();
        }
        await this.startServer();
    }
    getApiUrl() {
        const port = this.getConfig('serverPort', 8080);
        return `http://localhost:${port}/api/v1`;
    }
    isServerRunning() {
        return this.isRunning;
    }
    getServerStatus() {
        return {
            isRunning: this.isRunning,
            port: this.getConfig('serverPort', 8080),
            grpcPort: this.getConfig('grpcPort', 50051)
        };
    }
    async getDashboardData() {
        if (!this.isRunning) {
            throw new Error('Server is not running.');
        }
        const port = this.getConfig('serverPort', 8080);
        try {
            const response = await axios_1.default.get(`http://localhost:${port}/api/v1/stats`);
            return response.data;
        }
        catch (error) {
            if (axios_1.default.isAxiosError(error)) {
                throw new Error(`Failed to fetch dashboard data from server: ${error.message}`);
            }
            throw error;
        }
    }
    async getTools() {
        if (!this.isRunning) {
            return [];
        }
        const port = this.getConfig('serverPort', 8080);
        try {
            const response = await axios_1.default.get(`http://localhost:${port}/api/v1/tools`);
            return response.data.tools || [];
        }
        catch (error) {
            console.error('Failed to get tools:', error);
            return [];
        }
    }
    async getTool(toolName) {
        if (!this.isRunning) {
            return undefined;
        }
        const port = this.getConfig('serverPort', 8080);
        try {
            const response = await axios_1.default.get(`http://localhost:${port}/api/v1/tools/${toolName}`);
            return response.data;
        }
        catch (error) {
            console.error(`Failed to get tool ${toolName}:`, error);
            return undefined;
        }
    }
    async getAgents() {
        if (!this.isRunning) {
            return [];
        }
        const port = this.getConfig('serverPort', 8080);
        try {
            const response = await axios_1.default.get(`http://localhost:${port}/api/v1/agents`);
            return response.data.agents || [];
        }
        catch (error) {
            console.error('Failed to get agents:', error);
            return [];
        }
    }
    async executeTool(toolName, parameters) {
        if (!this.isRunning)
            throw new Error('Server is not running.');
        try {
            const response = await axios_1.default.post(`${this.getApiUrl()}/tools/${toolName}/execute`, { args: parameters });
            return response.data;
        }
        catch (error) {
            this.handleApiError(`Execution of tool '${toolName}' failed`, error);
            throw error;
        }
    }
    async getServerStats() {
        if (!this.isRunning)
            return null;
        try {
            const response = await axios_1.default.get(`${this.getApiUrl()}/stats`);
            return response.data;
        }
        catch (error) {
            this.handleApiError('Failed to fetch server stats', error);
            return null;
        }
    }
    handleApiError(message, error) {
        const errorMessage = axios_1.default.isAxiosError(error) && error.response?.data?.error
            ? error.response.data.error
            : error.message;
        vscode.window.showErrorMessage(`${message}: ${errorMessage}`);
        this.outputChannel.appendLine(`[API ERROR] ${message}: ${errorMessage}`);
    }
    showOutputChannel() {
        this.outputChannel.show();
    }
    dispose() {
        this.stopServerProcess();
        this.outputChannel.dispose();
        this.stateChangeEmitter.dispose();
    }
}
exports.ServerManager = ServerManager;
//# sourceMappingURL=serverManager.js.map