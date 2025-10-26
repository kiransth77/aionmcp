const fs = require('fs');
const path = require('path');

// Create bin directory in extension
const binDir = path.join(__dirname, '..', 'bin');
if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
}

// Build optimized binary for current platform only
console.log('Building optimized AionMCP binary for current platform...');

const { execSync } = require('child_process');
const parentDir = path.join(__dirname, '..', '..');

try {
    process.chdir(parentDir);
    const platform = process.platform;
    const extension = platform === 'win32' ? '.exe' : '';
    const outputPath = path.join('vscode-extension', 'bin', `aionmcp${extension}`);
    
    console.log(`Building optimized binary for ${platform}...`);
    
    // Build with size optimizations: -s removes symbol table, -w removes debug info
    execSync(`go build -ldflags="-s -w" -o ${outputPath} cmd/server/main.go`, { stdio: 'inherit' });
    
    // Make executable on Unix systems
    if (platform !== 'win32') {
        fs.chmodSync(path.join(__dirname, '..', 'bin', `aionmcp${extension}`), 0o755);
    }
    
    console.log(`✅ Built optimized binary: ${outputPath}`);
    
    // Show file size
    const binaryPath = path.join(__dirname, '..', 'bin', `aionmcp${extension}`);
    if (fs.existsSync(binaryPath)) {
        const stats = fs.statSync(binaryPath);
        const sizeMB = (stats.size / (1024 * 1024)).toFixed(2);
        console.log(`📦 Binary size: ${sizeMB} MB`);
    }
    
} catch (error) {
    console.error('❌ Failed to build AionMCP binary:', error.message);
    process.exit(1);
}

console.log('✅ Optimized binary build completed!');