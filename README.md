# ai-rss

A togo **AI data-source** plugin — fetch and parse **RSS / Atom / JSON** feeds so `ai-rag` ingest and agents can pull article content.

```
togo install togo-framework/ai-rss
```

## Use
- Go: `rss.FromKernel(k).Fetch(ctx, "https://blog.example.com/feed.xml")` → `[]Item{Title,Link,Content,Published}`
- REST: `POST /api/ai/rss` `{"url":"https://…/feed.xml"}`

Part of the [togo AI kit](https://to-go.dev/ai). MIT.
