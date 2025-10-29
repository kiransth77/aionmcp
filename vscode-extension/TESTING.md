# AionMCP Extension - Local Testing Guide

## 🧪 **Testing Methods**

### Method 1: Install VSIX Package (Recommended)
```bash
cd vscode-extension
code --install-extension aionmcp-extension-0.1.0.vsix
code ..  # Open workspace
```

### Method 2: Development Mode (F5 Debug)
1. Open `vscode-extension/` folder in VS Code
2. Press `F5` to launch Extension Development Host
3. Test in the new VS Code window

## ✅ **Testing Checklist**

### 1. Extension Activation
- [ ] AionMCP icon appears in Activity Bar (left sidebar)
- [ ] Extension loads without errors
- [ ] Status bar shows "AionMCP: Not Running"

### 2. Server Management
- [ ] Click AionMCP icon → Server Status view shows
- [ ] "Start AionMCP Server" button works
- [ ] Server starts and status changes to "Running"
- [ ] Port 8080 is accessible (server responds)
- [ ] Stop/Restart buttons work correctly

### 3. Tool Discovery
- [ ] Tools view populates with discovered tools
- [ ] Tools are categorized properly
- [ ] Refresh button works
- [ ] Search/filter functionality works

### 4. Tool Execution
- [ ] Right-click tool → "Execute Tool" works
- [ ] Tool Executor webview opens
- [ ] Parameter forms render correctly
- [ ] Tool execution returns results
- [ ] Results display properly formatted

### 5. API Spec Import
- [ ] "Import API Specification" button works
- [ ] Can select OpenAPI/GraphQL files
- [ ] New tools appear after import
- [ ] Error handling for invalid specs

### 6. Dashboard & Monitoring
- [ ] "Show Dashboard" command works
- [ ] Dashboard displays server stats
- [ ] Agent connections show (if any)
- [ ] Real-time updates work

### 7. Configuration
- [ ] Settings under File → Preferences → Settings → AionMCP
- [ ] Server path configuration works
- [ ] Port settings are respected
- [ ] Auto-start option works

## 🔍 **Manual Testing Steps**

### Step 1: Basic Functionality
1. **Open VS Code** with AionMCP workspace
2. **Look for AionMCP icon** in the left sidebar
3. **Click the icon** to open AionMCP views
4. **Check Server Status** view shows "Not Running"

### Step 2: Start Server
1. **Click "Start AionMCP Server"** in Server Status view
2. **Watch terminal output** for startup messages
3. **Verify status changes** to "Running" with green icon
4. **Check port accessibility**: Open browser to `http://localhost:8080`

### Step 3: Test Tool Discovery
1. **Check Tools view** populates with default tools
2. **Click refresh button** to reload tools
3. **Expand categories** to see tool organization
4. **Try search** if you have many tools

### Step 4: Execute a Tool
1. **Right-click any tool** in Tools view
2. **Select "Execute Tool"** or "Open Tool Executor"
3. **Fill in parameters** in the webview form
4. **Click "Execute"** and verify results appear
5. **Check error handling** with invalid parameters

### Step 5: Test API Import
1. **Click "Import API Specification"** in Tools view
2. **Select a spec file** from `examples/specs/`
3. **Verify new tools** appear in the list
4. **Test imported tools** work correctly

### Step 6: Dashboard Testing
1. **Open Command Palette** (`Ctrl+Shift+P`)
2. **Type "AionMCP: Show Dashboard"**
3. **Verify dashboard opens** with server metrics
4. **Check real-time updates** by executing tools

## 🐛 **Common Issues & Solutions**

### Extension Not Loading
- **Check VS Code Developer Tools**: `Help → Toggle Developer Tools`
- **Look for errors** in Console tab
- **Reload window**: `Ctrl+Shift+P` → "Developer: Reload Window"

### Server Won't Start
- **Check binary exists**: `vscode-extension/bin/aionmcp.exe`
- **Verify permissions**: Binary should be executable
- **Check port conflicts**: Port 8080 may be in use
- **Review logs**: AionMCP Server output channel

### Tools Not Loading
- **Server must be running** first
- **Check API connectivity**: `http://localhost:8080/api/tools`
- **Import API specs** if no tools are found
- **Refresh tools view** manually

### WebViews Not Opening
- **Check VS Code version**: Requires VS Code 1.85.0+
- **Disable other extensions** that might conflict
- **Clear VS Code cache**: Restart VS Code
- **Check console errors** in Developer Tools

## 📊 **Performance Testing**

### Package Size Verification
```bash
Get-Item aionmcp-extension-0.1.0.vsix | Select-Object Name, @{Name="Size(MB)"; Expression={[math]::Round($_.Length/1MB, 2)}}
```
**Expected**: ~9.4 MB (74% reduction from original 36.45 MB)

### Memory Usage
- **Monitor VS Code memory** with extension active
- **Check server process** memory usage
- **Verify no memory leaks** during extended use

### Startup Performance
- **Extension activation time**: Should be < 1 second
- **Server startup time**: Should be < 5 seconds
- **Tool discovery time**: Should be < 3 seconds

## ✅ **Success Criteria**

- [ ] Extension installs without errors
- [ ] All UI components render correctly
- [ ] Server starts and stops reliably
- [ ] Tools can be discovered and executed
- [ ] WebViews open and function properly
- [ ] No console errors or warnings
- [ ] Package size is optimized (~9.4 MB)
- [ ] Cross-platform binary detection works

## 🚀 **Next Steps After Testing**

1. **Document any issues found**
2. **Fix critical bugs** if discovered
3. **Update version** if changes made
4. **Re-package** if fixes applied
5. **Prepare for marketplace submission**

---

**🎯 Goal**: Verify the extension works flawlessly before merging to main branch!