# ai-rss — documentation

  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />

## Overview

Package rss is a togo AI data-source plugin: fetch + parse RSS/Atom feeds so
ai-rag ingest and agents can pull article content. Registers an "ai-rss"
service on the kernel and a REST endpoint: POST /api/ai/rss {"url":"…"}.

## Install

```bash
togo install togo-framework/ai-rss
```

A capability plugin — it self-registers on boot; no driver selector needed.

## Configuration

Environment variables read by this plugin (extracted from the source):

_No environment variables read directly (uses the kernel/base config or the app DB)._

## Usage

```go
// A data source for ai-rag / agents: fetch/scrape/search web content.
src := rss.FromKernel(k)
docs, err := src.Fetch(ctx, "https://example.com")
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/ai-rss
- README: ../README.md
