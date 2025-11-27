# Publishing AionMCP to the New MCP Registry

The new MCP Registry at https://registry.modelcontextprotocol.io is the official way to publish MCP servers. Follow these steps:

## Step 1: Install mcp-publisher CLI

```bash
# macOS with Homebrew (recommended)
brew install mcp-publisher

# Or download the binary
curl -L "https://github.com/modelcontextprotocol/registry/releases/latest/download/mcp-publisher_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/').tar.gz" | tar xz mcp-publisher && sudo mv mcp-publisher /usr/local/bin/

# Verify installation
mcp-publisher --help
```

## Step 2: Create server.json (Already Done! ✅)

The `server.json` file is already created in your project root:
- Location: `/Users/kiran/Documents/GitHub/aionmcp/server.json`
- Follows MCP Registry schema v2025-10-17
- Includes binary distribution for macOS (amd64, arm64), Linux (amd64, arm64), Windows (amd64)
- Includes Docker distribution

## Step 3: Build and Release Binaries (REQUIRED)

Before publishing, you need to build multi-platform binaries and create a GitHub release.

### 3a. Build Binaries

```bash
cd /Users/kiran/Documents/GitHub/aionmcp

# Create dist directory
mkdir -p dist

# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o dist/aionmcp-server ./cmd/server
tar czf dist/aionmcp-v0.1.0-darwin-amd64.tar.gz -C dist aionmcp-server
rm dist/aionmcp-server

# Build for macOS (ARM64/Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dist/aionmcp-server ./cmd/server
tar czf dist/aionmcp-v0.1.0-darwin-arm64.tar.gz -C dist aionmcp-server
rm dist/aionmcp-server

# Build for Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o dist/aionmcp-server ./cmd/server
tar czf dist/aionmcp-v0.1.0-linux-amd64.tar.gz -C dist aionmcp-server
rm dist/aionmcp-server

# Build for Linux (arm64)
GOOS=linux GOARCH=arm64 go build -o dist/aionmcp-server ./cmd/server
tar czf dist/aionmcp-v0.1.0-linux-arm64.tar.gz -C dist aionmcp-server
rm dist/aionmcp-server

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o dist/aionmcp-server.exe ./cmd/server
cd dist
zip aionmcp-v0.1.0-windows-amd64.zip aionmcp-server.exe
rm aionmcp-server.exe
cd ..
```

### 3b. Calculate SHA256 Hashes

```bash
cd dist
sha256sum aionmcp-v0.1.0-*.tar.gz aionmcp-v0.1.0-*.zip
```

Copy the hash values and update `server.json`:
- Replace `"hash": "sha256:PLACEHOLDER"` with actual SHA256 hashes

### 3c. Create GitHub Release

```bash
# Tag the release
git tag -a v0.1.0 -m "Release v0.1.0 - Model-independent MCP server"
git push origin v0.1.0

# Create release on GitHub with the binaries
# Or use GitHub CLI:
gh release create v0.1.0 dist/*.tar.gz dist/*.zip \
  --title "AionMCP v0.1.0" \
  --notes "Initial release of AionMCP - model-independent MCP server with OpenAPI, GraphQL, and AsyncAPI support"
```

## Step 4: Update server.json with Actual Hashes

Edit `/Users/kiran/Documents/GitHub/aionmcp/server.json` and replace `"hash": "sha256:PLACEHOLDER"` with real SHA256 values from Step 3b.

## Step 5: Authenticate with MCP Registry

```bash
cd /Users/kiran/Documents/GitHub/aionmcp

# Login with GitHub
mcp-publisher login github

# Follow the prompts:
# 1. Go to https://github.com/login/device
# 2. Enter the code shown in terminal
# 3. Authorize the application
# 4. Return to terminal - should see "Successfully authenticated!"
```

## Step 6: Publish to Registry

```bash
cd /Users/kiran/Documents/GitHub/aionmcp

# Publish your server
mcp-publisher publish

# You should see:
# ✓ Successfully published
# ✓ Server io.github.kiransth77/aionmcp version 0.1.0
```

## Step 7: Verify Publication

```bash
# Search for your server
curl "https://registry.modelcontextprotocol.io/v0.1/servers?search=io.github.kiransth77/aionmcp"

# Should return your server metadata in JSON
```

## Troubleshooting

| Error | Solution |
|-------|----------|
| "Invalid or expired Registry JWT token" | Re-authenticate: `mcp-publisher login github` |
| "You do not have permission to publish" | Your server name must start with `io.github.kiransth77/` (matches your GitHub username) |
| "Registry validation failed" | Ensure server.json follows the schema and all hashes are valid |
| "File not found" | Make sure binaries are in the correct URLs from server.json |

## Next Steps

1. ✅ Create `server.json` ← DONE
2. ⏳ Build multi-platform binaries
3. ⏳ Create GitHub release with binaries
4. ⏳ Update hashes in server.json
5. ⏳ Authenticate with registry
6. ⏳ Publish to registry

## Resources

- **Registry Quickstart**: https://github.com/modelcontextprotocol/registry/blob/main/docs/modelcontextprotocol-io/quickstart.mdx
- **Package Types**: https://github.com/modelcontextprotocol/registry/blob/main/docs/modelcontextprotocol-io/package-types.mdx
- **Authentication**: https://github.com/modelcontextprotocol/registry/blob/main/docs/modelcontextprotocol-io/authentication.mdx
- **Live Registry**: https://registry.modelcontextprotocol.io/docs
