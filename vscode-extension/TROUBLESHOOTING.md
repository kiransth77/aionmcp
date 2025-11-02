# AionMCP Extension - Quick Verification Guide

## 🔍 **Where to Find AionMCP Extension**

### 1. Activity Bar (Left Sidebar)
- Look for a **server/process icon** (📡) in the left sidebar
- Click it to open the AionMCP views
- If not visible, the extension may not be activated

### 2. Command Palette (Ctrl+Shift+P)
Type any of these commands:
- `AionMCP: Start Server`
- `AionMCP: Stop Server`
- `AionMCP: Show Dashboard`
- `AionMCP: Import API Specification`

### 3. VS Code Settings
- Go to **File → Preferences → Settings**
- Search for "AionMCP"
- You should see configuration options like:
  - Server Path
  - Server Port
  - Auto Start
  - Log Level

### 4. Extensions View
- Open **Extensions** (Ctrl+Shift+X)
- Search for "AionMCP"
- Should show "AionMCP - Autonomous MCP Server" as installed

## 🚨 **Troubleshooting**

### Extension Not Visible?

1. **Check if installed**:
   ```bash
   code --list-extensions | grep aionmcp
   ```

2. **Reload VS Code Window**:
   - `Ctrl+Shift+P` → "Developer: Reload Window"

3. **Check Developer Console**:
   - `Help → Toggle Developer Tools`
   - Look for errors in Console tab

4. **Force activation**:
   - `Ctrl+Shift+P` → Type "AionMCP" to trigger activation

### Extension Installed but Not Working?

1. **Check extension status**:
   - `Ctrl+Shift+P` → "Extensions: Show Running Extensions"
   - Look for AionMCP in the list

2. **Check activation events**:
   - Extension activates "onStartupFinished"
   - May take a few seconds after VS Code starts

3. **Manual activation**:
   - Try running any AionMCP command to force activation

## ✅ **Expected Behavior**

When working correctly, you should see:

### In Activity Bar:
- **AionMCP icon** (server/process symbol)
- Three views when clicked:
  - Tools
  - Connected Agents  
  - Server Status

### In Command Palette:
- 9 AionMCP commands available
- All starting with "AionMCP:"

### In Settings:
- "AionMCP" section with 6 configuration options

## 🎯 **Quick Test**

1. Open VS Code
2. Press `Ctrl+Shift+P`
3. Type "AionMCP: Start"
4. Should see "AionMCP: Start AionMCP Server" command
5. If yes → Extension is working! ✅
6. If no → Extension needs troubleshooting ❌

---

**Note**: AionMCP is NOT a "Language Tool" - it's a standalone MCP server extension, so it won't appear in language-specific configuration areas.