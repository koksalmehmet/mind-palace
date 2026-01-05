# Next.js 16 & React 19 Best Practices - 2025

**Research Date:** January 5, 2026  
**Target:** Mind Palace Documentation Site  
**Current Stack:** Next.js 16.1.1, React 19.2.3, Nextra 4.6.1

## Executive Summary

Your Mind Palace docs site is already using cutting-edge versions (Next.js 16.1.1 and React 19.2.3) with a solid foundation. This document provides 2025 best practices and actionable recommendations for optimization.

---

## 1. Next.js 16 App Router & Server Components Best Practices

### Current State ✅
- Using App Router with proper file-based routing
- Server components by default (no unnecessary 'use client' directives)
- Static export configured for GitHub Pages
- Proper async/await in server components

### 2025 Best Practices

#### Server Components (Default)
```tsx
// ✅ GOOD: Server Component (default, no directive needed)
export default async function Page(props: PageProps) {
  const params = await props.params
  const result = await importPage(params.mdxPath)
  return <Content {...result} />
}
```

#### Client Components (Only When Needed)
```tsx
// Only use 'use client' for:
// - Interactive features (onClick, useState)
// - Browser APIs (localStorage, window)
// - React hooks
'use client'

export function InteractiveSearch() {
  const [query, setQuery] = useState('')
  return <input onChange={(e) => setQuery(e.target.value)} />
}
```

#### New Next.js 16 Features to Adopt

**1. Enhanced Caching with `cacheLife`**
```tsx
import { unstable_cacheLife as cacheLife } from 'next/cache'

export async function getDocContent() {
  'use cache'
  cacheLife('hours')
  
  // This data is cached for hours
  return fetchMarkdownContent()
}
```

**2. Better Error Boundaries**
```tsx
// app/error.tsx
'use client'

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <div>
      <h2>Documentation Error</h2>
      <p>{error.message}</p>
      <button onClick={reset}>Try again</button>
    </div>
  )
}
```

**3. Loading UI Enhancements**
```tsx
// app/[[...mdxPath]]/loading.tsx
export default function Loading() {
  return (
    <div className="animate-pulse">
      <div className="h-8 bg-gray-200 rounded w-3/4 mb-4" />
      <div className="h-4 bg-gray-200 rounded w-full mb-2" />
      <div className="h-4 bg-gray-200 rounded w-5/6" />
    </div>
  )
}
```

### Recommendations for Mind Palace

1. **Add Error Boundary:** Create `apps/docs/app/error.tsx` for graceful error handling
2. **Add Loading State:** Create `apps/docs/app/[[...mdxPath]]/loading.tsx` for better UX
3. **Optimize Metadata:** Use `generateMetadata` more strategically with caching
4. **Consider Route Groups:** Organize docs by feature areas using `(groups)`

---

## 2. React 19 Features & Breaking Changes

### What's New in React 19

#### 1. Actions & Form Improvements
```tsx
// Modern form handling with React 19
'use client'

import { useActionState } from 'react'

export function FeedbackForm() {
  const [state, formAction] = useActionState(submitFeedback, null)
  
  return (
    <form action={formAction}>
      <input name="feedback" required />
      <button type="submit">Send Feedback</button>
      {state?.message && <p>{state.message}</p>}
    </form>
  )
}
```

#### 2. Document Metadata API
```tsx
// In any component (no more Head component needed)
export function DocsPage() {
  return (
    <>
      <title>Mind Palace | Documentation</title>
      <meta name="description" content="Learn about Mind Palace" />
      <link rel="canonical" href="https://example.com/docs" />
      <YourContent />
    </>
  )
}
```

#### 3. Asset Loading API
```tsx
// Preload resources programmatically
import { preload } from 'react-dom'

function MyComponent() {
  preload('/fonts/custom-font.woff2', { as: 'font', type: 'font/woff2' })
  return <div>Content</div>
}
```

#### 4. use() Hook for Promises
```tsx
'use client'

import { use } from 'react'

function DocContent({ contentPromise }) {
  // Unwrap promise directly in render
  const content = use(contentPromise)
  return <div>{content}</div>
}
```

### Breaking Changes (Already Handled)

✅ **Props as Promises** - Your code already uses `await props.params`
✅ **Refs as Props** - Standard ref forwarding works
✅ **Context API** - No changes needed for your use case

### Recommendations

1. **Add Feedback Form:** Use React 19 form actions for doc feedback
2. **Preload Critical Assets:** Add font/image preloading in layout
3. **Simplify Metadata:** Can use inline title/meta tags if needed

---

## 3. Image Optimization & Asset Loading

### Current State
```tsx
// next.config.mjs
images: {
  unoptimized: true,  // ⚠️ For static export
}
```

### 2025 Best Practices for Static Sites

#### Static Image Optimization
```tsx
import Image from 'next/image'

// ✅ Use Image component even with unoptimized: true
// Benefits: Lazy loading, proper sizing, better CLS
<Image
  src="/mind-palace/images/architecture.png"
  alt="Architecture diagram"
  width={800}
  height={600}
  loading="lazy"
  placeholder="blur"
  blurDataURL="data:image/svg+xml;base64,..."
/>
```

#### Asset Preloading Strategy
```tsx
// app/layout.tsx
export default function RootLayout({ children }) {
  return (
    <html>
      <head>
        {/* Critical CSS */}
        <link rel="preload" href="/styles/critical.css" as="style" />
        
        {/* Custom fonts */}
        <link
          rel="preload"
          href="/fonts/inter-var.woff2"
          as="font"
          type="font/woff2"
          crossOrigin="anonymous"
        />
        
        {/* Hero image */}
        <link rel="preload" href="/hero.webp" as="image" />
      </head>
      <body>{children}</body>
    </html>
  )
}
```

#### Modern Image Formats
```bash
# Convert images to WebP/AVIF for better performance
# Use sharp or squoosh for preprocessing

npm install -D sharp
```

```javascript
// scripts/optimize-images.js
const sharp = require('sharp')
const fs = require('fs').promises

async function optimizeImages() {
  const images = await fs.readdir('apps/docs/public/images')
  
  for (const img of images) {
    if (img.endsWith('.png') || img.endsWith('.jpg')) {
      await sharp(`public/images/${img}`)
        .webp({ quality: 85 })
        .toFile(`public/images/${img.replace(/\.(png|jpg)$/, '.webp')}`)
    }
  }
}
```

### Recommendations

1. **Create Image Optimization Script:** Convert PNGs to WebP
2. **Add Blur Placeholders:** Generate data URLs for critical images
3. **Implement Picture Element:** Serve WebP with PNG fallback
4. **Lazy Load Everything:** Ensure all non-critical images use lazy loading
5. **Add Logo Component:** Optimize Mind Palace logo delivery

---

## 4. Static Export & GitHub Pages Deployment

### Current Configuration ✅
```javascript
// next.config.mjs
export default withNextra({
  output: 'export',  // ✅ Static export enabled
  images: {
    unoptimized: true,  // ✅ Required for static export
  },
  basePath: process.env.NODE_ENV === 'production' ? '/mind-palace' : '',  // ✅ Good
})
```

### 2025 Best Practices

#### Enhanced Build Configuration
```javascript
// next.config.mjs
import nextra from 'nextra'

const withNextra = nextra({
  latex: true,
  defaultShowCopyCode: true,
  // Add these for better static export
  theme: 'nextra-theme-docs',
  themeConfig: './theme.config.tsx'
})

export default withNextra({
  output: 'export',
  distDir: 'out',
  
  // Improved image handling
  images: {
    unoptimized: true,
    remotePatterns: [], // No remote images for static
  },
  
  // Base path handling
  basePath: process.env.NODE_ENV === 'production' ? '/mind-palace' : '',
  assetPrefix: process.env.NODE_ENV === 'production' ? '/mind-palace' : '',
  
  // Trailing slashes for GitHub Pages
  trailingSlash: true,
  
  // Skip type checking during build (do it separately)
  typescript: {
    ignoreBuildErrors: false,
  },
  
  // Optimize CSS
  experimental: {
    optimizeCss: true, // Requires critters
  },
  
  // Better compression
  compress: true,
  
  // Generate sitemap
  generateBuildId: async () => {
    return 'mind-palace-' + new Date().toISOString()
  },
})
```

#### Optimized GitHub Pages Deployment
```yaml
# .github/workflows/deploy-docs.yml
name: Deploy Docs

on:
  push:
    branches: [main]
    paths:
      - 'apps/docs/**'
      - '.github/workflows/deploy-docs.yml'

# Allow only one concurrent deployment
concurrency:
  group: pages
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: apps/docs/package-lock.json
      
      - name: Install dependencies
        run: cd apps/docs && npm ci
      
      - name: Build docs
        run: cd apps/docs && npm run build
        env:
          NODE_ENV: production
      
      - name: Add .nojekyll
        run: touch apps/docs/out/.nojekyll
      
      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: apps/docs/out

  deploy:
    needs: build
    runs-on: ubuntu-latest
    
    permissions:
      pages: write
      id-token: write
    
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

#### Custom 404 Page
```tsx
// app/not-found.tsx
import Link from 'next/link'

export default function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <h1 className="text-6xl font-bold">404</h1>
      <h2 className="text-2xl mt-4">Page Not Found</h2>
      <p className="mt-2 text-gray-600">
        The documentation page you're looking for doesn't exist.
      </p>
      <Link 
        href="/"
        className="mt-8 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
      >
        Return to Docs Home
      </Link>
    </div>
  )
}
```

### Recommendations

1. **Add GitHub Actions Workflow:** Automate deployment
2. **Create Custom 404:** Better UX for missing pages
3. **Add Sitemap Generation:** Help search engines crawl
4. **Implement trailing slashes:** Better GitHub Pages compatibility
5. **Add robots.txt:** `apps/docs/public/robots.txt`
6. **Create .nojekyll:** Prevent Jekyll processing

---

## 5. TypeScript Strict Mode in Next.js

### Current State ✅
Your `tsconfig.json` is already excellent:
```json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "strictBindCallApply": true,
    "strictPropertyInitialization": true,
    "noImplicitThis": true,
    "alwaysStrict": true
  }
}
```

### 2025 Additional Recommendations

#### Enhanced TypeScript Configuration
```json
{
  "compilerOptions": {
    // Your existing strict options ✅
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "strictBindCallApply": true,
    "strictPropertyInitialization": true,
    "noImplicitThis": true,
    "alwaysStrict": true,
    
    // Additional 2025 best practices
    "noUncheckedIndexedAccess": true,  // Safer array access
    "noUnusedLocals": true,             // Catch unused variables
    "noUnusedParameters": true,         // Catch unused params
    "noFallthroughCasesInSwitch": true, // Safer switch statements
    "noImplicitReturns": true,          // All code paths return
    "forceConsistentCasingInFileNames": true,
    
    // Modern module resolution
    "moduleResolution": "bundler",
    "module": "esnext",
    "target": "ES2022",  // Update from ES2017
    
    // Path aliases for cleaner imports
    "baseUrl": ".",
    "paths": {
      "@/*": ["./*"],
      "@/components/*": ["./components/*"],
      "@/content/*": ["./content/*"]
    }
  }
}
```

#### Type-Safe Route Parameters
```typescript
// types/routes.ts
export type DocRoute = 
  | ['getting-started']
  | ['getting-started', 'concepts']
  | ['getting-started', 'workflows']
  | ['features', 'brain']
  | ['features', 'corridors']
  // ... etc

export type PageProps = {
  params: Promise<{
    mdxPath?: DocRoute
  }>
}
```

#### Strict Component Props
```typescript
// Define all prop types strictly
type LayoutProps = Readonly<{
  children: React.ReactNode
  params?: {
    locale?: string
  }
}>

export default function Layout({ children, params }: LayoutProps) {
  return <div>{children}</div>
}
```

### Recommendations

1. **Update tsconfig.json:** Add the additional strict checks
2. **Define Route Types:** Type-safe route parameters
3. **Update target:** ES2022 for better performance
4. **Add Path Aliases:** Cleaner imports with @/ prefix
5. **Type MDX Components:** Ensure all custom components are typed

---

## 6. Link Checking Automation

### Problem
Documentation often has broken internal/external links that degrade UX.

### 2025 Solution: Automated Link Checking

#### Option 1: markdown-link-check (GitHub Actions)
```yaml
# .github/workflows/check-links.yml
name: Check Documentation Links

on:
  push:
    branches: [main]
    paths:
      - 'apps/docs/content/**/*.mdx'
  pull_request:
    paths:
      - 'apps/docs/content/**/*.mdx'
  schedule:
    - cron: '0 0 * * 1'  # Weekly on Monday

jobs:
  check-links:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Check MDX links
        uses: gaurav-nelson/github-action-markdown-link-check@v1
        with:
          use-quiet-mode: 'yes'
          use-verbose-mode: 'yes'
          config-file: '.markdown-link-check.json'
          folder-path: 'apps/docs/content'
```

#### Configuration File
```json
// .markdown-link-check.json
{
  "ignorePatterns": [
    {
      "pattern": "^http://localhost"
    }
  ],
  "replacementPatterns": [
    {
      "pattern": "^/",
      "replacement": "http://localhost:3000/"
    }
  ],
  "httpHeaders": [
    {
      "urls": ["https://github.com"],
      "headers": {
        "Accept": "text/html"
      }
    }
  ],
  "timeout": "20s",
  "retryOn429": true,
  "retryCount": 3,
  "aliveStatusCodes": [200, 206, 301, 302]
}
```

#### Option 2: Playwright for Deployed Site Testing
```typescript
// tests/link-checker.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Documentation Links', () => {
  const baseUrl = 'https://koksalmehmet.github.io/mind-palace'
  
  const pages = [
    '/getting-started',
    '/getting-started/concepts',
    '/features/brain',
    '/reference/cli',
  ]
  
  for (const page of pages) {
    test(`check links on ${page}`, async ({ page: pwPage }) => {
      await pwPage.goto(baseUrl + page)
      
      // Get all links
      const links = await pwPage.locator('a[href]').all()
      
      for (const link of links) {
        const href = await link.getAttribute('href')
        if (!href || href.startsWith('#')) continue
        
        // Check internal links
        if (href.startsWith('/')) {
          const response = await pwPage.request.get(baseUrl + href)
          expect(response.ok()).toBeTruthy()
        }
        
        // Check external links (sample)
        if (href.startsWith('http')) {
          const response = await pwPage.request.get(href)
          expect(response.status()).toBeLessThan(400)
        }
      }
    })
  }
})
```

#### Option 3: Custom Script for Pre-commit Hook
```javascript
// scripts/check-doc-links.js
const fs = require('fs').promises
const path = require('path')
const https = require('https')

async function checkDocLinks() {
  const contentDir = 'apps/docs/content'
  const files = await getMarkdownFiles(contentDir)
  
  let hasErrors = false
  
  for (const file of files) {
    const content = await fs.readFile(file, 'utf-8')
    const links = extractLinks(content)
    
    for (const link of links) {
      if (link.startsWith('/')) {
        // Check internal link
        const exists = await checkInternalLink(link)
        if (!exists) {
          console.error(`❌ Broken internal link in ${file}: ${link}`)
          hasErrors = true
        }
      } else if (link.startsWith('http')) {
        // Check external link
        const ok = await checkExternalLink(link)
        if (!ok) {
          console.warn(`⚠️  External link issue in ${file}: ${link}`)
        }
      }
    }
  }
  
  if (hasErrors) {
    process.exit(1)
  }
}

function extractLinks(content) {
  const linkRegex = /\[([^\]]+)\]\(([^)]+)\)/g
  const links = []
  let match
  
  while ((match = linkRegex.exec(content)) !== null) {
    links.push(match[2])
  }
  
  return links
}

async function checkInternalLink(link) {
  // Convert /getting-started to content/getting-started/index.mdx
  const possiblePaths = [
    `apps/docs/content${link}.mdx`,
    `apps/docs/content${link}/index.mdx`,
  ]
  
  for (const p of possiblePaths) {
    try {
      await fs.access(p)
      return true
    } catch {}
  }
  
  return false
}

checkDocLinks()
```

#### Pre-commit Hook
```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Checking documentation links..."
node scripts/check-doc-links.js

if [ $? -ne 0 ]; then
    echo "❌ Link check failed. Please fix broken links before committing."
    exit 1
fi
```

### Recommendations

1. **Add GitHub Action:** Weekly automated link checking
2. **Create Configuration:** `.markdown-link-check.json`
3. **Add npm Script:** `"check-links": "node scripts/check-doc-links.js"`
4. **Optional Pre-commit Hook:** Catch issues before commit
5. **Monitor External Links:** Schedule weekly checks

---

## 7. MDX Integration & Performance

### Current State ✅
- Using Nextra 4.6.1 (modern MDX solution)
- Proper MDX component integration
- Server-side rendering

### 2025 MDX Best Practices

#### Custom MDX Components
```tsx
// mdx-components.tsx
import { useMDXComponents as getDocsMDXComponents } from 'nextra-theme-docs'
import Image from 'next/image'
import { Callout } from './components/Callout'
import { CodeBlock } from './components/CodeBlock'

const docsComponents = getDocsMDXComponents()

export function useMDXComponents(components?: Record<string, React.ComponentType>) {
  return {
    ...docsComponents,
    
    // Enhanced image component
    img: (props: any) => (
      <Image
        {...props}
        alt={props.alt || ''}
        width={800}
        height={600}
        className="rounded-lg"
        loading="lazy"
      />
    ),
    
    // Custom callout boxes
    Callout,
    
    // Syntax highlighted code blocks
    pre: CodeBlock,
    
    // External link with icon
    a: ({ href, children, ...props }: any) => {
      const isExternal = href?.startsWith('http')
      return (
        <a
          href={href}
          {...props}
          {...(isExternal && {
            target: '_blank',
            rel: 'noopener noreferrer',
          })}
        >
          {children}
          {isExternal && <span className="ml-1">↗</span>}
        </a>
      )
    },
    
    ...components,
  }
}
```

#### Custom Components for Docs

**Callout Component**
```tsx
// components/Callout.tsx
import { ReactNode } from 'react'

type CalloutType = 'info' | 'warning' | 'error' | 'success'

export function Callout({
  type = 'info',
  children,
}: {
  type?: CalloutType
  children: ReactNode
}) {
  const styles = {
    info: 'bg-blue-50 border-blue-200 text-blue-900',
    warning: 'bg-yellow-50 border-yellow-200 text-yellow-900',
    error: 'bg-red-50 border-red-200 text-red-900',
    success: 'bg-green-50 border-green-200 text-green-900',
  }
  
  const icons = {
    info: 'ℹ️',
    warning: '⚠️',
    error: '❌',
    success: '✅',
  }
  
  return (
    <div className={`border-l-4 p-4 my-4 rounded ${styles[type]}`}>
      <span className="mr-2">{icons[type]}</span>
      {children}
    </div>
  )
}
```

Usage in MDX:
```mdx
# Getting Started

<Callout type="info">
  Mind Palace requires Go 1.23 or later.
</Callout>

<Callout type="warning">
  Alpha version - Breaking changes expected
</Callout>
```

#### MDX Performance Optimization

**1. Code Splitting**
```tsx
// Lazy load heavy components
import dynamic from 'next/dynamic'

const InteractiveDiagram = dynamic(() => import('./InteractiveDiagram'), {
  loading: () => <p>Loading diagram...</p>,
  ssr: false, // Client-side only
})
```

**2. Content Caching**
```typescript
// lib/mdx-cache.ts
import { cache } from 'react'
import { importPage } from 'nextra/pages'

// React cache() ensures single fetch per request
export const getCachedPage = cache(async (path?: string[]) => {
  return importPage(path)
})
```

**3. Optimized Frontmatter**
```mdx
---
title: "Mind Palace CLI Reference"
description: "Complete guide to Mind Palace CLI commands"
lastUpdated: "2026-01-05"
keywords: ["cli", "commands", "reference"]
---

# CLI Reference
...
```

#### Static Content Generation
```typescript
// Optimize static params generation
export async function generateStaticParams() {
  const paths = await getAllDocPaths()
  
  return paths.map((path) => ({
    mdxPath: path.split('/').filter(Boolean),
  }))
}

async function getAllDocPaths() {
  // Recursively find all MDX files
  const contentDir = 'content'
  // ... implementation
}
```

### Recommendations

1. **Add Custom MDX Components:**
   - Callout boxes for tips/warnings
   - Code blocks with copy button
   - Image component with lazy loading
   - External link indicator

2. **Implement Content Caching:** Use React `cache()` for repeated queries

3. **Add Frontmatter Validation:** Ensure all docs have required metadata

4. **Create Component Library:** Reusable doc components in `/components`

5. **Optimize Bundle:** 
   - Dynamic imports for heavy components
   - Remove unused Nextra features

6. **Add Search:** Integrate Algolia DocSearch or Pagefind for static search

---

## Priority Implementation Roadmap

### Phase 1: Quick Wins (1-2 hours)
1. ✅ Add error.tsx for better error handling
2. ✅ Add loading.tsx for loading states
3. ✅ Add not-found.tsx for 404 page
4. ✅ Update tsconfig with stricter checks
5. ✅ Add .nojekyll to output

### Phase 2: Enhanced UX (2-4 hours)
1. 🔄 Create custom MDX components (Callout, CodeBlock)
2. 🔄 Implement link checking automation
3. 🔄 Add GitHub Actions for deployment
4. 🔄 Optimize images (convert to WebP)
5. 🔄 Add sitemap generation

### Phase 3: Advanced Features (4-8 hours)
1. ⏳ Implement static search (Pagefind)
2. ⏳ Add dark mode toggle
3. ⏳ Create interactive diagrams
4. ⏳ Add doc versioning
5. ⏳ Implement feedback system

### Phase 4: Performance & Analytics (2-4 hours)
1. ⏳ Add Web Vitals tracking
2. ⏳ Implement content preloading
3. ⏳ Optimize bundle size
4. ⏳ Add performance monitoring
5. ⏳ Create lighthouse CI

---

## Monitoring & Maintenance

### Performance Metrics to Track
```javascript
// app/layout.tsx
import { SpeedInsights } from '@vercel/speed-insights/next'
import { Analytics } from '@vercel/analytics/react'

export default function RootLayout({ children }) {
  return (
    <html>
      <body>
        {children}
        <SpeedInsights />
        <Analytics />
      </body>
    </html>
  )
}
```

### Automated Testing
```yaml
# .github/workflows/test-docs.yml
name: Test Documentation

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      
      - name: Install
        run: cd apps/docs && npm ci
      
      - name: Lint
        run: cd apps/docs && npm run lint
      
      - name: Type check
        run: cd apps/docs && npx tsc --noEmit
      
      - name: Build
        run: cd apps/docs && npm run build
      
      - name: Check links
        run: npm run check-links
```

---

## Key Takeaways

### What You're Doing Right ✅
1. Using latest Next.js 16 and React 19
2. Proper TypeScript strict mode configuration
3. Static export for GitHub Pages
4. Server components by default
5. Good project structure

### Quick Wins to Implement 🚀
1. Add error/loading/not-found pages
2. Implement link checking automation
3. Create custom MDX components
4. Add GitHub Actions for deployment
5. Optimize images with WebP

### Future Enhancements 🔮
1. Static search with Pagefind
2. Interactive documentation features
3. Doc versioning system
4. Performance monitoring
5. User feedback system

---

## Additional Resources

- [Next.js 16 Documentation](https://nextjs.org/docs)
- [React 19 Release Notes](https://react.dev/blog/2024/12/05/react-19)
- [Nextra Documentation](https://nextra.site)
- [MDX Documentation](https://mdxjs.com)
- [GitHub Pages Custom Domains](https://docs.github.com/pages)
- [Web.dev Performance](https://web.dev/learn-performance/)

---

**Next Steps:** Pick items from Phase 1 and implement them first. They provide immediate UX improvements with minimal effort.
