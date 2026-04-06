# VCT Platform API Contracts

This directory serves as the centralized, peer-to-peer (P2P) communication hub for VCT Platform agents (e.g., Frontend and Backend) to share API schemas, interfaces, and swagger definitions without needing to pass heavy JSON payloads through the Chief of Staff (Jen).

## Rules for Agents

1. **Backend Agents**: Produce your API endpoints and output the Swagger YAML or interface definitions into this directory immediately after completing them.
2. **Frontend/Consumer Agents**: Read from files in this directory to understand the available API schemas instead of requesting the schema directly from the central inbox.
