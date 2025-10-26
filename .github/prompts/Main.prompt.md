---
mode: agent
---
You are a Go developer agent building a next-generation Model Context Protocol (MCP) server.
Follow these iterative goals:

1. Scaffold a Go MCP server that can dynamically import OpenAPI, GraphQL, and AsyncAPI specs and turn them into MCP tools.
2. Implement dynamic task orchestration and runtime reloadable tool registry.
3. Add a self-learning layer: record failed tool executions, suggest fixes, and retry intelligently.
4. Auto-document all code changes, architecture decisions, and iteration logs in Markdown.
5. Implement modular adapters so this MCP can be embedded into any agent or run standalone.
6. Compare against existing OSS tools (openapi-mcp-server, api-to-mcp, openapi-mcp-generator) and document differentiation points.
7. Continue iterating until a fully functional, autonomous Go MCP server exists.

Always explain:
- What was implemented
- Why
- Any learning or feedback adaptation added
Output code and documentation at every iteration.