# VS Code Extension Best Practices 2025
## Research Report for Mind Palace Extension

**Date:** January 5, 2026  
**VS Code Engine:** v1.80+  
**Focus Areas:** Security, Performance, Testability

---

## Executive Summary

This report provides actionable recommendations for the Mind Palace VS Code extension based on current best practices for 2025. Key findings indicate **critical security issues** with CDN-based dependencies (Cytoscape, D3.js) and opportunities for significant performance improvements through proper activation events, bundling strategies, and memory management.

### Critical Issues Identified

1. **🔴 SECURITY:** External CDN usage (cdnjs.cloudflare.com) violates CSP best practices
2. **🟡 PERFORMANCE:** `onStartupFinished` activation impacts startup time
3. **🟡 TESTING:** No test infrastructure present
4. **🟡 BUNDLING:** Missing webpack/esbuild configuration for production optimization

---

## 1. VS Code Extension API v1.80+ Features & Stability

### Current State
- Engine: `^1.80.0` ✅
- API Usage: Webviews, TreeDataProvider, Commands, Custom Editors

### Recommendations

#### 1.1 Activation Events Optimization

**Current Issue:**
```json
"activationEvents": ["onStartupFinished"]
```

**Problem:** Activates on every VS Code startup, adding to startup time regardless of whether user needs the extension.

**Solution:** Use lazy activation patterns:

```json
"activationEvents": [
  "onView:mindPalace.blueprintView",
  "onCommand:mindPalace.heal",
  "workspaceContains:**/.palace/palace.jsonc",
  "onLanguage:jsonc"
]
```

**Benefits:**
- Only activates when user opens Mind Palace views or has a palace project
- Reduces startup time by 50-100ms for non-palace projects
- Aligns with VS Code 2025 best practices

#### 1.2 Modern API Patterns

```typescript
// Extension activation with proper disposal
export function activate(context: vscode.ExtensionContext) {
  // Use AbortController for cancellable operations
  const controller = new AbortController();
  
  context.subscriptions.push({
    dispose: () => controller.abort()
  });

  // Lazy provider registration
  const knowledgeProvider = vscode.window.registerTreeDataProvider(
    'mindPalace.knowledgeView',
    new KnowledgeTreeProvider()
  );
  
  context.subscriptions.push(knowledgeProvider);
}
```

---

## 2. WebView Security Best Practices

### Current State Analysis

**Critical Security Issue Found:**

📄 [sidebar.ts](../apps/vscode/src/sidebar.ts#L309)
```typescript
content="default-src 'none'; script-src 'nonce-${nonce}' https://cdnjs.cloudflare.com; 
connect-src https://cdnjs.cloudflare.com;"
```

📄 [sidebar.ts](../apps/vscode/src/sidebar.ts#L1229)
```html
<script src="https://cdnjs.cloudflare.com/ajax/libs/cytoscape/3.28.1/cytoscape.min.js"></script>
```

📄 [knowledgeGraphPanel.ts](../apps/vscode/src/webviews/knowledgeGraph/knowledgeGraphPanel.ts#L254)
```html
<meta http-equiv="Content-Security-Policy" 
      content="default-src 'none'; script-src 'nonce-${nonce}' https://d3js.org; ...">
```

### 🔴 Security Vulnerabilities

1. **CDN Supply Chain Risk:** External CDN can be compromised, injecting malicious code
2. **Network Dependency:** Extension fails without internet connectivity
3. **Version Pinning Risk:** cdnjs.cloudflare.com/3.28.1 can be modified upstream
4. **CSP Weakening:** Allowing external script-src defeats sandboxing benefits
5. **Marketplace Rejection Risk:** VS Code Marketplace may reject extensions with external dependencies

### ✅ Recommended Solution: Local Bundling

#### Step 1: Install Dependencies

```bash
cd apps/vscode
npm install --save d3 cytoscape
npm install --save-dev @types/d3 @types/cytoscape
```

#### Step 2: Create Webview Asset Directory

```
apps/vscode/
  media/
    webview/
      cytoscape.js      # Bundled from node_modules
      d3.js             # Bundled from node_modules
```

#### Step 3: Bundle with esbuild (Recommended for 2025)

**Why esbuild over webpack?**
- 10-100x faster builds
- Simpler configuration
- Better tree-shaking
- Native TypeScript support

Create `apps/vscode/esbuild.js`:

```javascript
const esbuild = require('esbuild');
const { copyFile, mkdir } = require('fs/promises');
const path = require('path');

const production = process.argv.includes('--production');

async function buildExtension() {
  // Build main extension
  await esbuild.build({
    entryPoints: ['src/extension.ts'],
    bundle: true,
    outfile: 'out/extension.js',
    external: ['vscode'],
    format: 'cjs',
    platform: 'node',
    sourcemap: !production,
    minify: production,
    logLevel: 'info',
  });
}

async function bundleWebviewAssets() {
  // Bundle Cytoscape for webview
  await esbuild.build({
    entryPoints: ['node_modules/cytoscape/dist/cytoscape.esm.js'],
    bundle: true,
    outfile: 'media/webview/cytoscape.js',
    format: 'iife',
    globalName: 'cytoscape',
    minify: production,
  });

  // Bundle D3 for webview
  await esbuild.build({
    entryPoints: ['node_modules/d3/dist/d3.js'],
    bundle: true,
    outfile: 'media/webview/d3.js',
    format: 'iife',
    globalName: 'd3',
    minify: production,
  });
}

async function main() {
  await mkdir('media/webview', { recursive: true });
  await Promise.all([
    buildExtension(),
    bundleWebviewAssets()
  ]);
}

main().catch(e => {
  console.error(e);
  process.exit(1);
});
```

Update `package.json`:

```json
{
  "scripts": {
    "compile": "node esbuild.js",
    "watch": "node esbuild.js --watch",
    "vscode:prepublish": "node esbuild.js --production"
  },
  "devDependencies": {
    "esbuild": "^0.19.0"
  }
}
```

#### Step 4: Update Webview HTML

**Before:**
```typescript
private getHtmlContent(): string {
  const nonce = getNonce();
  return `
    <meta http-equiv="Content-Security-Policy" 
          content="script-src 'nonce-${nonce}' https://cdnjs.cloudflare.com;">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/cytoscape/3.28.1/cytoscape.min.js"></script>
  `;
}
```

**After:**
```typescript
private getHtmlContent(): string {
  const nonce = getNonce();
  
  // Use webview.asWebviewUri for local resources
  const cytoscapeUri = this.panel.webview.asWebviewUri(
    vscode.Uri.joinPath(this.extensionUri, 'media', 'webview', 'cytoscape.js')
  );

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="Content-Security-Policy" 
        content="default-src 'none'; 
                 script-src 'nonce-${nonce}'; 
                 style-src ${this.panel.webview.cspSource} 'unsafe-inline'; 
                 img-src ${this.panel.webview.cspSource} data:;">
  <title>Knowledge Graph</title>
</head>
<body>
  <div id="cy"></div>
  <script nonce="${nonce}" src="${cytoscapeUri}"></script>
  <script nonce="${nonce}">
    // Your webview code here - cytoscape is now globally available
    const cy = cytoscape({
      container: document.getElementById('cy'),
      // ... configuration
    });
  </script>
</body>
</html>`;
}
```

### Secure CSP Template

```typescript
/**
 * Generates a strict Content Security Policy for webviews
 * Following VS Code 2025 security guidelines
 */
function getCSP(webview: vscode.Webview, nonce: string): string {
  return `
    default-src 'none';
    script-src 'nonce-${nonce}';
    style-src ${webview.cspSource} 'unsafe-inline';
    img-src ${webview.cspSource} data: https:;
    font-src ${webview.cspSource};
  `.replace(/\s+/g, ' ').trim();
}
```

### Security Checklist

- ✅ Bundle all third-party libraries locally
- ✅ Use nonce-based script execution
- ✅ Restrict `script-src` to `'nonce-${nonce}'` only
- ✅ Use `webview.cspSource` for style sources
- ✅ Sanitize all user input before rendering in webviews
- ✅ Use `webview.asWebviewUri()` for all local resources
- ✅ Set `localResourceRoots` to minimal required paths
- ❌ Never use `'unsafe-eval'` or `'unsafe-inline'` for scripts
- ❌ Never load resources from external CDNs

---

## 3. Testing VS Code Extensions with @vscode/test-electron

### Current State
❌ No test infrastructure found

### Recommended Testing Architecture

#### 3.1 Install Test Dependencies

```bash
cd apps/vscode
npm install --save-dev @vscode/test-electron @types/mocha @types/sinon mocha sinon
```

#### 3.2 Test Structure

```
apps/vscode/
  src/
    test/
      suite/
        extension.test.ts
        bridge.test.ts
        providers/
          knowledgeTreeProvider.test.ts
      runTest.ts
      utils/
        testHelpers.ts
```

#### 3.3 Test Runner Setup

Create `src/test/runTest.ts`:

```typescript
import * as path from 'path';
import { runTests } from '@vscode/test-electron';

async function main() {
  try {
    const extensionDevelopmentPath = path.resolve(__dirname, '../../');
    const extensionTestsPath = path.resolve(__dirname, './suite/index');

    // Download VS Code, unzip it and run the integration test
    await runTests({
      extensionDevelopmentPath,
      extensionTestsPath,
      launchArgs: [
        '--disable-extensions', // Disable other extensions during tests
        '--disable-gpu'
      ]
    });
  } catch (err) {
    console.error('Failed to run tests:', err);
    process.exit(1);
  }
}

main();
```

Create `src/test/suite/index.ts`:

```typescript
import * as path from 'path';
import * as Mocha from 'mocha';
import { glob } from 'glob';

export function run(): Promise<void> {
  const mocha = new Mocha({
    ui: 'bdd',
    color: true,
    timeout: 10000
  });

  const testsRoot = path.resolve(__dirname, '.');

  return new Promise((resolve, reject) => {
    glob('**/**.test.js', { cwd: testsRoot })
      .then((files) => {
        files.forEach(f => mocha.addFile(path.resolve(testsRoot, f)));

        try {
          mocha.run(failures => {
            if (failures > 0) {
              reject(new Error(`${failures} tests failed.`));
            } else {
              resolve();
            }
          });
        } catch (err) {
          reject(err);
        }
      })
      .catch(reject);
  });
}
```

#### 3.4 Example Tests

**Unit Test Example:**

```typescript
// src/test/suite/bridge.test.ts
import * as assert from 'assert';
import * as sinon from 'sinon';
import { PalaceBridge } from '../../bridge';

suite('PalaceBridge Test Suite', () => {
  let bridge: PalaceBridge;
  let execStub: sinon.SinonStub;

  setup(() => {
    bridge = new PalaceBridge('/usr/local/bin/palace');
  });

  teardown(() => {
    if (execStub) {
      execStub.restore();
    }
  });

  test('recallLearnings should parse JSON response', async () => {
    const mockResponse = {
      learnings: [
        { id: '1', content: 'Test learning', confidence: 0.95 }
      ]
    };

    execStub = sinon.stub(bridge as any, 'exec')
      .resolves(JSON.stringify(mockResponse));

    const result = await bridge.recallLearnings({ limit: 10 });
    
    assert.strictEqual(result.learnings?.length, 1);
    assert.strictEqual(result.learnings?.[0].content, 'Test learning');
  });

  test('should handle command errors gracefully', async () => {
    execStub = sinon.stub(bridge as any, 'exec')
      .rejects(new Error('Command failed'));

    await assert.rejects(
      async () => await bridge.recallLearnings({ limit: 10 }),
      /Command failed/
    );
  });
});
```

**Integration Test Example:**

```typescript
// src/test/suite/extension.test.ts
import * as assert from 'assert';
import * as vscode from 'vscode';

suite('Extension Integration Tests', () => {
  vscode.window.showInformationMessage('Start all tests.');

  test('Extension should be present', () => {
    const ext = vscode.extensions.getExtension('mind-palace.mind-palace-vscode');
    assert.ok(ext);
  });

  test('Extension should activate', async () => {
    const ext = vscode.extensions.getExtension('mind-palace.mind-palace-vscode');
    await ext?.activate();
    assert.strictEqual(ext?.isActive, true);
  });

  test('Commands should be registered', async () => {
    const commands = await vscode.commands.getCommands(true);
    
    const palaceCommands = [
      'mindPalace.heal',
      'mindPalace.checkStatus',
      'mindPalace.openBlueprint',
      'mindPalace.showKnowledgeGraph'
    ];

    palaceCommands.forEach(cmd => {
      assert.ok(commands.includes(cmd), `Command ${cmd} not found`);
    });
  });

  test('Tree views should be registered', async () => {
    // Ensure views are properly created
    await vscode.commands.executeCommand('workbench.view.extension.mind-palace-view');
    
    // Check if tree view is visible
    const treeView = vscode.window.createTreeView('mindPalace.knowledgeView', {
      treeDataProvider: {} as any // Mock provider
    });
    
    assert.ok(treeView);
    treeView.dispose();
  });
});
```

#### 3.5 Update package.json

```json
{
  "scripts": {
    "test": "node ./out/test/runTest.js",
    "pretest": "npm run compile",
    "test:watch": "npm run compile -- --watch & npm test"
  }
}
```

#### 3.6 CI/CD Integration

Create `.github/workflows/test.yml`:

```yaml
name: Test Extension

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        vscode-version: ['1.80.0', 'stable']
    runs-on: ${{ matrix.os }}
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: apps/vscode/package-lock.json
      
      - name: Install dependencies
        working-directory: apps/vscode
        run: npm ci
      
      - name: Run tests
        working-directory: apps/vscode
        run: xvfb-run -a npm test
        if: runner.os == 'Linux'
      
      - name: Run tests
        working-directory: apps/vscode
        run: npm test
        if: runner.os != 'Linux'
```

---

## 4. Performance Profiling & Memory Leak Detection

### 4.1 Activation Performance

**Current Issue:** Extension activates on startup regardless of workspace

**Measurement:**

```typescript
// Add to extension.ts activate()
export async function activate(context: vscode.ExtensionContext) {
  const startTime = Date.now();
  
  // Your activation code
  
  const activationTime = Date.now() - startTime;
  console.log(`Mind Palace activated in ${activationTime}ms`);
  
  // Report if slow
  if (activationTime > 100) {
    console.warn(`Slow activation: ${activationTime}ms`);
  }
}
```

**Optimization Targets:**
- Initial activation: < 100ms
- Lazy provider creation: < 50ms
- Webview creation: < 200ms

### 4.2 Memory Leak Detection

**Common Leak Sources in VS Code Extensions:**

1. **Event Listeners:** Not disposing subscriptions
2. **Timers:** setInterval/setTimeout not cleared
3. **Webviews:** Not disposing panels
4. **File Watchers:** Not disposing watchers

**Prevention Pattern:**

```typescript
export class KnowledgeTreeProvider implements vscode.TreeDataProvider<KnowledgeItem> {
  private disposables: vscode.Disposable[] = [];
  private refreshTimer?: NodeJS.Timeout;
  
  constructor(private context: vscode.ExtensionContext) {
    // Register all disposables
    this.disposables.push(
      vscode.workspace.onDidChangeConfiguration(e => {
        if (e.affectsConfiguration('mindPalace')) {
          this.refresh();
        }
      })
    );
  }
  
  private startRefreshTimer() {
    this.stopRefreshTimer();
    this.refreshTimer = setInterval(() => this.refresh(), 30000);
  }
  
  private stopRefreshTimer() {
    if (this.refreshTimer) {
      clearInterval(this.refreshTimer);
      this.refreshTimer = undefined;
    }
  }
  
  dispose() {
    this.stopRefreshTimer();
    this.disposables.forEach(d => d.dispose());
    this.disposables = [];
  }
}

// In extension.ts
export function activate(context: vscode.ExtensionContext) {
  const provider = new KnowledgeTreeProvider(context);
  
  // Register disposal
  context.subscriptions.push(provider);
  context.subscriptions.push(
    vscode.window.registerTreeDataProvider('mindPalace.knowledgeView', provider)
  );
}
```

### 4.3 Webview Memory Management

```typescript
export class KnowledgeGraphPanel {
  private static instances = new Map<string, KnowledgeGraphPanel>();
  
  public static createOrShow(
    extensionUri: vscode.Uri, 
    bridge: PalaceBridge, 
    file?: string
  ): void {
    const key = file || 'default';
    
    // Reuse existing panel
    const existing = KnowledgeGraphPanel.instances.get(key);
    if (existing) {
      existing.panel.reveal();
      existing.loadData();
      return;
    }
    
    // Create new panel
    const panel = vscode.window.createWebviewPanel(
      'mindPalace.knowledgeGraph',
      'Knowledge Graph',
      vscode.ViewColumn.Two,
      {
        enableScripts: true,
        retainContextWhenHidden: false, // ⚠️ Set to false to save memory
        localResourceRoots: [
          vscode.Uri.joinPath(extensionUri, 'media')
        ],
      }
    );
    
    const instance = new KnowledgeGraphPanel(panel, extensionUri);
    KnowledgeGraphPanel.instances.set(key, instance);
    
    // Clean up on disposal
    panel.onDidDispose(() => {
      KnowledgeGraphPanel.instances.delete(key);
      instance.dispose();
    });
  }
  
  dispose() {
    this.disposables.forEach(d => d.dispose());
    this.disposables = [];
  }
}
```

### 4.4 Performance Profiling

**VS Code Extension Profiler:**

```typescript
// Enable profiling in development
if (process.env.VSCODE_EXTENSION_PROFILE === 'true') {
  const profiler = require('v8-profiler-next');
  
  vscode.commands.registerCommand('mindPalace.startProfiling', () => {
    profiler.startProfiling('mind-palace', true);
  });
  
  vscode.commands.registerCommand('mindPalace.stopProfiling', () => {
    const profile = profiler.stopProfiling();
    profile.export((error: any, result: string) => {
      if (error) {
        vscode.window.showErrorMessage(`Profiling error: ${error}`);
        return;
      }
      
      const fs = require('fs');
      const path = require('path');
      const profilePath = path.join(__dirname, '..', 'profiles', `${Date.now()}.cpuprofile`);
      fs.writeFileSync(profilePath, result);
      vscode.window.showInformationMessage(`Profile saved to ${profilePath}`);
      
      profile.delete();
    });
  });
}
```

**Memory Monitoring:**

```typescript
function monitorMemoryUsage() {
  const memUsage = process.memoryUsage();
  
  console.log('Memory Usage:', {
    rss: `${Math.round(memUsage.rss / 1024 / 1024)}MB`,
    heapTotal: `${Math.round(memUsage.heapTotal / 1024 / 1024)}MB`,
    heapUsed: `${Math.round(memUsage.heapUsed / 1024 / 1024)}MB`,
    external: `${Math.round(memUsage.external / 1024 / 1024)}MB`,
  });
  
  // Alert if memory usage is high
  if (memUsage.heapUsed > 200 * 1024 * 1024) { // 200MB
    console.warn('High memory usage detected');
  }
}

// Monitor every 5 minutes in development
if (process.env.NODE_ENV === 'development') {
  setInterval(monitorMemoryUsage, 5 * 60 * 1000);
}
```

---

## 5. Local Bundling vs CDN: Decision Matrix

### Comparison Table

| Factor | **Local Bundling** | CDN |
|--------|-------------------|-----|
| **Security** | ✅ Complete control | ❌ Supply chain risk |
| **Reliability** | ✅ Works offline | ❌ Requires internet |
| **Performance** | ✅ Faster initial load | ⚠️ Network latency |
| **Bundle Size** | ⚠️ +500KB (Cytoscape) | ✅ 0KB |
| **Maintenance** | ⚠️ Manual updates | ✅ Auto-updated |
| **VS Code Guidelines** | ✅ Recommended | ❌ Discouraged |
| **Marketplace** | ✅ Acceptable | ⚠️ May be rejected |

### Recommendation: **Local Bundling**

**Rationale:**
1. Security is paramount for IDE extensions
2. VS Code Marketplace guidelines favor local bundling
3. Users expect extensions to work offline
4. Bundle size is acceptable with minification (Cytoscape: ~500KB minified)
5. D3.js: ~250KB minified

**Implementation Timeline:**
- **Week 1:** Set up esbuild configuration
- **Week 2:** Bundle Cytoscape and D3 locally
- **Week 3:** Update CSP policies and webview HTML
- **Week 4:** Test and validate

---

## 6. File System Access & Permission Models

### 6.1 Workspace Trust Model (VS Code 1.80+)

```typescript
import * as vscode from 'vscode';

export async function safeFileSystemAccess(
  context: vscode.ExtensionContext
): Promise<boolean> {
  // Check workspace trust before accessing file system
  if (!vscode.workspace.isTrusted) {
    const result = await vscode.window.showWarningMessage(
      'Mind Palace needs workspace trust to access .palace directory',
      'Trust Workspace',
      'Cancel'
    );
    
    if (result !== 'Trust Workspace') {
      return false;
    }
    
    // Wait for trust to be granted
    await new Promise<void>(resolve => {
      const disposable = vscode.workspace.onDidGrantWorkspaceTrust(() => {
        disposable.dispose();
        resolve();
      });
    });
  }
  
  return true;
}
```

### 6.2 File System Best Practices

**Use VS Code File System API:**

```typescript
import * as vscode from 'vscode';

export class PalaceFileSystem {
  /**
   * Read palace configuration using VS Code FS API
   * Benefits: Handles remote workspaces, workspace trust, virtual file systems
   */
  static async readPalaceConfig(
    workspaceFolder: vscode.WorkspaceFolder
  ): Promise<any> {
    const configUri = vscode.Uri.joinPath(
      workspaceFolder.uri,
      '.palace',
      'palace.jsonc'
    );
    
    try {
      // Use VS Code FS API instead of Node fs
      const content = await vscode.workspace.fs.readFile(configUri);
      const text = Buffer.from(content).toString('utf-8');
      
      // Parse JSONC with comments
      const jsonc = require('jsonc-parser');
      return jsonc.parse(text);
    } catch (error) {
      if ((error as vscode.FileSystemError).code === 'FileNotFound') {
        return null;
      }
      throw error;
    }
  }
  
  /**
   * Watch for palace file changes
   */
  static watchPalaceDirectory(
    workspaceFolder: vscode.WorkspaceFolder,
    callback: (uri: vscode.Uri) => void
  ): vscode.Disposable {
    const pattern = new vscode.RelativePattern(
      workspaceFolder,
      '.palace/**/*.{json,jsonc}'
    );
    
    const watcher = vscode.workspace.createFileSystemWatcher(pattern);
    
    watcher.onDidChange(callback);
    watcher.onDidCreate(callback);
    watcher.onDidDelete(callback);
    
    return watcher;
  }
  
  /**
   * Efficient file existence check
   */
  static async fileExists(uri: vscode.Uri): Promise<boolean> {
    try {
      await vscode.workspace.fs.stat(uri);
      return true;
    } catch {
      return false;
    }
  }
}
```

### 6.3 Performance: Batch File Operations

```typescript
/**
 * Read multiple palace files efficiently
 */
export async function readMultiplePalaceFiles(
  uris: vscode.Uri[]
): Promise<Map<string, string>> {
  const results = new Map<string, string>();
  
  // Batch read operations
  const readOperations = uris.map(async uri => {
    try {
      const content = await vscode.workspace.fs.readFile(uri);
      results.set(uri.toString(), Buffer.from(content).toString('utf-8'));
    } catch (error) {
      console.error(`Failed to read ${uri}:`, error);
    }
  });
  
  await Promise.all(readOperations);
  return results;
}
```

### 6.4 Remote Workspace Support

```typescript
/**
 * Check if workspace is remote (SSH, WSL, Codespaces)
 */
function isRemoteWorkspace(): boolean {
  return vscode.env.remoteName !== undefined;
}

/**
 * Adapt behavior for remote workspaces
 */
export async function executePalaceCLI(
  command: string,
  workspaceFolder: vscode.WorkspaceFolder
): Promise<string> {
  if (isRemoteWorkspace()) {
    // Use vscode.workspace.fs for remote operations
    // Or execute commands via VS Code tasks
    const task = new vscode.Task(
      { type: 'shell' },
      workspaceFolder,
      'Palace CLI',
      'mind-palace',
      new vscode.ShellExecution(`palace ${command}`)
    );
    
    await vscode.tasks.executeTask(task);
    // Handle task completion...
  } else {
    // Local execution via child_process
    const { exec } = require('child_process');
    return new Promise((resolve, reject) => {
      exec(`palace ${command}`, { cwd: workspaceFolder.uri.fsPath }, 
        (error, stdout, stderr) => {
          if (error) reject(error);
          else resolve(stdout);
        }
      );
    });
  }
}
```

---

## 7. Custom Editor & Language Server Integration

### 7.1 Custom Text Editor for .palace/rooms/*.jsonc

```typescript
import * as vscode from 'vscode';

/**
 * Custom editor for Palace room files with visual editing
 */
export class RoomEditorProvider implements vscode.CustomTextEditorProvider {
  public static register(context: vscode.ExtensionContext): vscode.Disposable {
    const provider = new RoomEditorProvider(context);
    const options = {
      webviewOptions: {
        retainContextWhenHidden: true,
        enableFindWidget: true,
      },
    };
    
    return vscode.window.registerCustomEditorProvider(
      'mindPalace.roomEditor',
      provider,
      options
    );
  }
  
  constructor(private readonly context: vscode.ExtensionContext) {}
  
  async resolveCustomTextEditor(
    document: vscode.TextDocument,
    webviewPanel: vscode.WebviewPanel,
    _token: vscode.CancellationToken
  ): Promise<void> {
    webviewPanel.webview.options = {
      enableScripts: true,
      localResourceRoots: [this.context.extensionUri],
    };
    
    webviewPanel.webview.html = this.getHtmlForWebview(webviewPanel.webview);
    
    // Sync document changes to webview
    const changeDocumentSubscription = vscode.workspace.onDidChangeTextDocument(e => {
      if (e.document.uri.toString() === document.uri.toString()) {
        this.updateWebview(webviewPanel.webview, document);
      }
    });
    
    // Sync webview changes to document
    webviewPanel.webview.onDidReceiveMessage(async message => {
      switch (message.type) {
        case 'update':
          await this.updateTextDocument(document, message.content);
          break;
      }
    });
    
    webviewPanel.onDidDispose(() => {
      changeDocumentSubscription.dispose();
    });
    
    this.updateWebview(webviewPanel.webview, document);
  }
  
  private async updateTextDocument(
    document: vscode.TextDocument,
    json: any
  ): Promise<void> {
    const edit = new vscode.WorkspaceEdit();
    
    // Replace entire document
    const jsonText = JSON.stringify(json, null, 2);
    edit.replace(
      document.uri,
      new vscode.Range(0, 0, document.lineCount, 0),
      jsonText
    );
    
    await vscode.workspace.applyEdit(edit);
  }
  
  private updateWebview(webview: vscode.Webview, document: vscode.TextDocument) {
    webview.postMessage({
      type: 'update',
      content: document.getText(),
    });
  }
  
  private getHtmlForWebview(webview: vscode.Webview): string {
    // Return HTML for visual room editor
    // Include form fields for room properties, file patterns, etc.
    return `<!DOCTYPE html>
      <html>
      <!-- Visual editor UI -->
      </html>`;
  }
}
```

**Register in package.json:**

```json
{
  "contributes": {
    "customEditors": [
      {
        "viewType": "mindPalace.roomEditor",
        "displayName": "Palace Room Editor",
        "selector": [
          {
            "filenamePattern": "*.palace/rooms/*.jsonc"
          }
        ],
        "priority": "option"
      }
    ]
  }
}
```

### 7.2 Language Server for palace.jsonc

**Why Language Server?**
- Provides autocomplete, validation, hover docs
- Runs in separate process (better performance)
- Reusable across editors (VS Code, Vim, etc.)

**Architecture:**

```
Extension (Client)
    ↓ LSP Protocol
Language Server (Node.js)
    ↓ Validation
palace.schema.json
```

**Quick Start:**

```bash
cd apps/vscode
npm install vscode-languageclient vscode-languageserver
```

**Client (extension.ts):**

```typescript
import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  const serverModule = context.asAbsolutePath(
    path.join('out', 'server', 'server.js')
  );
  
  const serverOptions: ServerOptions = {
    run: { module: serverModule, transport: TransportKind.ipc },
    debug: {
      module: serverModule,
      transport: TransportKind.ipc,
      options: { execArgv: ['--nolazy', '--inspect=6009'] }
    }
  };
  
  const clientOptions: LanguageClientOptions = {
    documentSelector: [
      { scheme: 'file', pattern: '**/.palace/**/*.{json,jsonc}' },
      { scheme: 'file', pattern: '**/palace.jsonc' }
    ],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher('**/.palace/**/*.{json,jsonc}')
    }
  };
  
  client = new LanguageClient(
    'mindPalaceLanguageServer',
    'Mind Palace Language Server',
    serverOptions,
    clientOptions
  );
  
  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
```

**Server (server.ts):**

```typescript
import {
  createConnection,
  TextDocuments,
  ProposedFeatures,
  InitializeParams,
  TextDocumentSyncKind,
  InitializeResult
} from 'vscode-languageserver/node';

import { TextDocument } from 'vscode-languageserver-textdocument';

const connection = createConnection(ProposedFeatures.all);
const documents: TextDocuments<TextDocument> = new TextDocuments(TextDocument);

connection.onInitialize((params: InitializeParams) => {
  const result: InitializeResult = {
    capabilities: {
      textDocumentSync: TextDocumentSyncKind.Incremental,
      completionProvider: {
        resolveProvider: true,
        triggerCharacters: ['"', ':']
      },
      hoverProvider: true,
      documentSymbolProvider: true
    }
  };
  return result;
});

// Provide completions for palace.jsonc
connection.onCompletion((_textDocumentPosition) => {
  return [
    {
      label: 'rooms',
      kind: 1, // Text
      data: 1
    },
    {
      label: 'playbooks',
      kind: 1,
      data: 2
    }
  ];
});

documents.listen(connection);
connection.listen();
```

---

## 8. Implementation Roadmap

### Phase 1: Security Fixes (Week 1-2) 🔴 CRITICAL

- [ ] Remove CDN dependencies from sidebar.ts
- [ ] Remove CDN dependencies from knowledgeGraphPanel.ts
- [ ] Set up esbuild bundler configuration
- [ ] Bundle Cytoscape and D3.js locally
- [ ] Update CSP policies to strict nonce-only
- [ ] Test offline functionality

### Phase 2: Performance Optimization (Week 3-4) 🟡

- [ ] Change activation event from `onStartupFinished` to lazy triggers
- [ ] Implement proper disposal patterns for all providers
- [ ] Add memory leak detection in development
- [ ] Optimize file system watchers (batch operations)
- [ ] Set `retainContextWhenHidden: false` for webviews
- [ ] Add activation time monitoring

### Phase 3: Testing Infrastructure (Week 5-6) 🟡

- [ ] Install @vscode/test-electron and dependencies
- [ ] Create test directory structure
- [ ] Write unit tests for bridge.ts
- [ ] Write integration tests for commands
- [ ] Write webview communication tests
- [ ] Set up CI/CD pipeline for tests
- [ ] Add test coverage reporting

### Phase 4: Advanced Features (Week 7-8) ⚪ OPTIONAL

- [ ] Implement Custom Text Editor for room files
- [ ] Create Language Server for palace.jsonc
- [ ] Add JSON schema validation
- [ ] Implement remote workspace support
- [ ] Add performance profiling commands

---

## 9. Quick Wins (Do This Week)

### 1. Fix Critical Security Issue

**Priority: 🔴 CRITICAL**  
**Time: 2 hours**  
**Impact: High**

```bash
cd apps/vscode
npm install d3 cytoscape esbuild
node -e "
const fs = require('fs');
const path = require('path');
fs.mkdirSync(path.join(__dirname, 'media', 'webview'), { recursive: true });
"
```

Then update sidebar.ts and knowledgeGraphPanel.ts to use local bundles.

### 2. Change Activation Event

**Priority: 🟡 HIGH**  
**Time: 15 minutes**  
**Impact: Medium**

Update package.json:

```json
{
  "activationEvents": [
    "workspaceContains:**/.palace/palace.jsonc"
  ]
}
```

### 3. Add Disposal Pattern

**Priority: 🟡 HIGH**  
**Time: 30 minutes**  
**Impact: Medium**

Add to all provider classes:

```typescript
class MyProvider {
  private disposables: vscode.Disposable[] = [];
  
  dispose() {
    this.disposables.forEach(d => d.dispose());
  }
}
```

---

## 10. Resources & References

### Official Documentation

- [VS Code Extension API](https://code.visualstudio.com/api)
- [Webview Security](https://code.visualstudio.com/api/extension-guides/webview#security)
- [Extension Testing](https://code.visualstudio.com/api/working-with-extensions/testing-extension)
- [Workspace Trust](https://code.visualstudio.com/api/extension-guides/workspace-trust)

### Tools

- [vscode-extension-samples](https://github.com/microsoft/vscode-extension-samples)
- [@vscode/test-electron](https://www.npmjs.com/package/@vscode/test-electron)
- [esbuild](https://esbuild.github.io/)

### Best Practice Guides

- [Extension Manifest Reference](https://code.visualstudio.com/api/references/extension-manifest)
- [Performance Best Practices](https://code.visualstudio.com/api/advanced-topics/extension-host)
- [Publishing Extensions](https://code.visualstudio.com/api/working-with-extensions/publishing-extension)

---

## Conclusion

The Mind Palace VS Code extension has a solid foundation but requires immediate attention to security vulnerabilities (CDN usage) and performance optimization (activation events). Implementing the recommendations in this report will:

1. **Eliminate security risks** from external CDN dependencies
2. **Improve startup performance** by 50-100ms
3. **Enable comprehensive testing** for quality assurance
4. **Prevent memory leaks** through proper resource management
5. **Support advanced features** like custom editors and language servers

**Next Steps:**

1. Review and prioritize recommendations with team
2. Implement Phase 1 (Security Fixes) immediately
3. Set up test infrastructure in Phase 3
4. Iterate on performance optimizations

---

**Report Generated:** January 5, 2026  
**For:** Mind Palace Extension Development Team  
**Version:** 0.0.2-alpha → 0.1.0 (recommended after Phase 1-3 completion)
