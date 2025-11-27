# AionMCP MCP Registry Submission Instructions

## Step 1: Fork the Repository

1. Go to https://github.com/modelcontextprotocol/servers
2. Click the **Fork** button in the top-right corner
3. This will create a fork at `https://github.com/kiransth77/servers`

## Step 2: Clone Your Fork

```bash
cd /tmp/mcp-registry/servers
# Add your fork as origin (if not already set)
git remote set-url origin https://github.com/kiransth77/servers.git
git remote add upstream https://github.com/modelcontextprotocol/servers.git
```

## Step 3: Create a Feature Branch

```bash
cd /tmp/mcp-registry/servers
git checkout -b feat/add-aionmcp
```

## Step 4: Add AionMCP Files

The files are already created at:
- `/tmp/mcp-registry/servers/src/aionmcp/server.json`
- `/tmp/mcp-registry/servers/src/aionmcp/README.md`

Verify they exist:

```bash
ls -la /tmp/mcp-registry/servers/src/aionmcp/
```

## Step 5: Commit and Push

```bash
cd /tmp/mcp-registry/servers

# Stage the new files
git add src/aionmcp/

# Commit with descriptive message
git commit -m "feat: add AionMCP MCP server

- Multi-format API specification support (OpenAPI, GraphQL, AsyncAPI)
- Model-independent REST API for any LLM integration
- Dynamic tool registration from API specs
- Self-learning engine with feedback collection
- VS Code extension with full IDE integration
- Hot-reload capability for spec updates"

# Push to your fork
git push origin feat/add-aionmcp
```

## Step 6: Create Pull Request

1. Go to your fork: https://github.com/kiransth77/servers
2. GitHub will show a prompt to create a PR from `feat/add-aionmcp`
3. Click **"Compare & pull request"**

4. Fill in the PR details:

**Title:**
```
feat: add AionMCP MCP server
```

**Description:**
```markdown
## AionMCP

Add AionMCP as a new MCP server to the registry.

### What is AionMCP?

AionMCP is a model-independent MCP server that dynamically imports OpenAPI, GraphQL, and AsyncAPI specifications and exposes them as tools to agents. It works with GitHub Copilot, Claude, or any LLM through a universal REST API.

### Key Features

- **Multi-format Support**: OpenAPI 3.0, GraphQL, AsyncAPI 2.0+
- **Model-Independent**: Universal HTTP REST API (not MCP-dependent)
- **Self-Learning**: Autonomous feedback collection and reflection engine
- **Hot-Reloadable**: Dynamically import and register new API specifications
- **VS Code Extension**: Full IDE integration with tool management UI
- **Auto-Documentation**: Generates changelog and reflection reports

### Quick Links

- Repository: https://github.com/kiransth77/aionmcp
- Documentation: https://github.com/kiransth77/aionmcp/blob/main/docs/GITHUB_COPILOT_INTEGRATION.md
- Setup Guide: https://github.com/kiransth77/aionmcp/blob/main/docs/SETUP_GITHUB_COPILOT.md

### Installation

```bash
# Download from releases
./aionmcp-server

# Or build from source
git clone https://github.com/kiransth77/aionmcp.git
cd aionmcp
go build -o aionmcp-server ./cmd/server
./aionmcp-server
```

### Use Cases

1. **Agent Integration** - Use with custom AI agents to perform API calls
2. **Multi-API Orchestration** - Chain operations across multiple APIs
3. **Development Tools** - Quick API testing and integration
4. **LLM Integration** - Enable language models to interact with APIs

### Files Changed

- Added `src/aionmcp/server.json` - Server metadata and configuration
- Added `src/aionmcp/README.md` - Installation and usage guide
```

## Step 7: Wait for Review

The MCP registry maintainers will review your PR. They typically check:
- ✅ Server functionality and usability
- ✅ Documentation completeness
- ✅ Compliance with registry standards
- ✅ No security concerns
- ✅ Unique value proposition

Review usually takes 1-2 weeks.

## Step 8: Address Feedback

If maintainers request changes:
1. Make the changes in your local branch
2. Commit with `git commit -m "feedback: address review comments"`
3. Push: `git push origin feat/add-aionmcp`
4. The PR will update automatically

## Step 9: Merge

Once approved, the maintainers will merge your PR. AionMCP will then be officially listed in the MCP registry!

---

## Quick Reference

**Files to commit:**
```
src/aionmcp/server.json
src/aionmcp/README.md
```

**Commit command:**
```bash
cd /tmp/mcp-registry/servers
git add src/aionmcp/
git commit -m "feat: add AionMCP MCP server"
git push origin feat/add-aionmcp
```

**PR URL (after pushing):**
https://github.com/kiransth77/servers/pull/new/feat/add-aionmcp

---

## Support

If you encounter any issues:
1. Check the [MCP Registry Contributing Guide](https://github.com/modelcontextprotocol/servers/blob/main/CONTRIBUTING.md)
2. Review other server examples in `src/` directory
3. Open an issue in the [MCP Registry repository](https://github.com/modelcontextprotocol/servers/issues)
