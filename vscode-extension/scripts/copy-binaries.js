const fs = require('fs');
const path = require('path');

// Create bin directory in extension
const binDir = path.join(__dirname, '..', 'bin');
if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
}

// Source binary paths (from parent project)
const parentBinDir = path.join(__dirname, '..', '..', 'bin');
const platforms = {
    'win32': 'aionmcp.exe',
    'linux': 'aionmcp',
    'darwin': 'aionmcp'
};

console.log('Copying AionMCP binaries...');

// Check if parent binaries exist
if (!fs.existsSync(parentBinDir)) {
    console.warn('Warning: Parent bin directory not found. Building AionMCP binary...');
    
    // Build the binary
    const { execSync } = require('child_process');
    const parentDir = path.join(__dirname, '..', '..');
    
    try {
        // Build for current platform
        process.chdir(parentDir);
        const platform = process.platform;
        const extension = platform === 'win32' ? '.exe' : '';
        const outputPath = path.join('bin', `aionmcp${extension}`);
        
        console.log(`Building for ${platform}...`);
        execSync(`go build -ldflags "-s -w" -o ${outputPath} cmd/server/main.go`, { stdio: 'inherit' });
        
        // Copy to extension bin directory
        const sourcePath = path.join(parentDir, outputPath);
        const targetPath = path.join(binDir, `aionmcp${extension}`);
        
        if (fs.existsSync(sourcePath)) {
            fs.copyFileSync(sourcePath, targetPath);
            fs.chmodSync(targetPath, 0o755); // Make executable
            console.log(`✅ Copied binary: ${targetPath}`);
        } else {
            console.error(`❌ Binary not found: ${sourcePath}`);
            process.exit(1);
        }
        
    } catch (error) {
        console.error('❌ Failed to build AionMCP binary:', error.message);
        process.exit(1);
    }
} else {
    // Copy existing binaries
    for (const [platform, filename] of Object.entries(platforms)) {
        const sourcePath = path.join(parentBinDir, filename);
        const targetPath = path.join(binDir, filename);
        
        if (fs.existsSync(sourcePath)) {
            fs.copyFileSync(sourcePath, targetPath);
            if (platform !== 'win32') {
                fs.chmodSync(targetPath, 0o755); // Make executable on Unix systems
            }
            console.log(`✅ Copied ${platform} binary: ${targetPath}`);
        } else {
            console.log(`⚠️  Binary not found for ${platform}: ${sourcePath}`);
        }
    }
}

console.log('✅ Binary copy process completed!');