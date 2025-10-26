#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

// Platform configurations
const platforms = {
    'win32': { 
        binary: 'aionmcp.exe',
        name: 'windows',
        goos: 'windows',
        goarch: 'amd64'
    },
    'linux': { 
        binary: 'aionmcp',
        name: 'linux',
        goos: 'linux',
        goarch: 'amd64'
    },
    'darwin': { 
        binary: 'aionmcp',
        name: 'macos',
        goos: 'darwin',
        goarch: 'amd64'
    }
};

async function buildForPlatform(platformKey, config) {
    console.log(`\n🔨 Building optimized package for ${config.name}...`);
    
    const binDir = path.join(__dirname, 'bin');
    const parentDir = path.join(__dirname, '..');
    
    // Clean bin directory
    if (fs.existsSync(binDir)) {
        fs.rmSync(binDir, { recursive: true, force: true });
    }
    fs.mkdirSync(binDir, { recursive: true });
    
    try {
        // Build optimized binary for target platform
        const outputPath = path.join('vscode-extension', 'bin', config.binary);
        const buildCmd = `go build -ldflags="-s -w" -o ${outputPath} cmd/server/main.go`;
        const env = {
            ...process.env,
            GOOS: config.goos,
            GOARCH: config.goarch,
            CGO_ENABLED: '0'
        };
        
        console.log(`Building: ${buildCmd}`);
        process.chdir(parentDir);
        execSync(buildCmd, { stdio: 'inherit', env });
        
        // Make executable on Unix
        if (platformKey !== 'win32') {
            const binaryPath = path.join(__dirname, 'bin', config.binary);
            fs.chmodSync(binaryPath, 0o755);
        }
        
        // Update .vscodeignore for this platform
        const vscodeignoreContent = `# Development files
.vscode/**
.vscode-test/**
src/**
.gitignore
.yarnrc
vsc-extension-quickstart.md
**/tsconfig.json
**/.eslintrc.json
**/.eslintrc.js
scripts/**
.eslintcache
*.vsix

# Source maps and TypeScript
**/*.map
**/*.ts
!out/**/*.js

# Dependencies
node_modules/**

# Platform-specific binaries (${config.name} only)
bin/**
!bin/${config.binary}

# Documentation (minimal for package size)
package-lock.json
CHANGELOG.md

# Media/resources
media/**
resources/**`;

        fs.writeFileSync(path.join(__dirname, '.vscodeignore'), vscodeignoreContent);
        
        // Build package
        process.chdir(__dirname);
        
        // Update package.json version for platform-specific build
        const packageJson = JSON.parse(fs.readFileSync('package.json', 'utf8'));
        const originalName = packageJson.name;
        const originalDisplayName = packageJson.displayName;
        
        packageJson.name = `${originalName}-${config.name}`;
        packageJson.displayName = `${originalDisplayName} (${config.name.toUpperCase()})`;
        
        fs.writeFileSync('package.json', JSON.stringify(packageJson, null, 2));
        
        // Package for this platform
        execSync('npx vsce package', { stdio: 'inherit' });
        
        // Restore original package.json
        packageJson.name = originalName;
        packageJson.displayName = originalDisplayName;
        fs.writeFileSync('package.json', JSON.stringify(packageJson, null, 2));
        
        // Show file size
        const vsixPattern = `aionmcp-extension-${config.name}-0.1.0.vsix`;
        const vsixFiles = fs.readdirSync('.').filter(f => f.includes(config.name) && f.endsWith('.vsix'));
        
        if (vsixFiles.length > 0) {
            const stats = fs.statSync(vsixFiles[0]);
            const sizeMB = (stats.size / (1024 * 1024)).toFixed(2);
            console.log(`✅ ${config.name} package: ${vsixFiles[0]} (${sizeMB} MB)`);
        }
        
    } catch (error) {
        console.error(`❌ Failed to build ${config.name} package:`, error.message);
        return false;
    }
    
    return true;
}

async function main() {
    console.log('🚀 Building optimized multi-platform packages...\n');
    
    const targetPlatforms = process.argv.slice(2);
    const buildPlatforms = targetPlatforms.length > 0 
        ? targetPlatforms 
        : Object.keys(platforms);
    
    let successCount = 0;
    
    for (const platformKey of buildPlatforms) {
        if (!platforms[platformKey]) {
            console.error(`❌ Unknown platform: ${platformKey}`);
            continue;
        }
        
        const success = await buildForPlatform(platformKey, platforms[platformKey]);
        if (success) successCount++;
    }
    
    console.log(`\n🎉 Successfully built ${successCount}/${buildPlatforms.length} platform packages!`);
    
    // Show all created packages
    const vsixFiles = fs.readdirSync('.').filter(f => f.endsWith('.vsix'));
    if (vsixFiles.length > 0) {
        console.log('\n📦 Created packages:');
        vsixFiles.forEach(file => {
            const stats = fs.statSync(file);
            const sizeMB = (stats.size / (1024 * 1024)).toFixed(2);
            console.log(`  • ${file} (${sizeMB} MB)`);
        });
    }
}

// Usage information
if (process.argv.includes('--help') || process.argv.includes('-h')) {
    console.log(`
Usage: node package-multi-platform.js [platforms...]

Platforms: win32, linux, darwin
Examples:
  node package-multi-platform.js           # Build all platforms
  node package-multi-platform.js win32     # Build Windows only
  node package-multi-platform.js linux darwin  # Build Linux and macOS

This script creates optimized, platform-specific VSIX packages with:
- Stripped binaries (no debug symbols)
- Single platform binary per package
- Minimal file inclusion for size optimization
`);
    process.exit(0);
}

main().catch(console.error);