"""
OpenCode Python Client SDK
Async client for interacting with OpenCode Pro Headless Daemon.
"""

import httpx
from typing import Dict, Any, Optional

class OpenCodeClient:
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url.rstrip("/")

    async def health(self) -> Dict[str, Any]:
        """Check OpenCode daemon health and introspected stack."""
        async with httpx.AsyncClient(timeout=10.0) as client:
            resp = await client.get(f"{self.base_url}/api/v1/health")
            resp.raise_for_status()
            return resp.json()

    async def optimize(self, prompt: str) -> Dict[str, Any]:
        """Optimize a raw prompt into a structured Golden Prompt."""
        async with httpx.AsyncClient(timeout=30.0) as client:
            resp = await client.post(
                f"{self.base_url}/api/v1/optimize",
                json={"prompt": prompt}
            )
            resp.raise_for_status()
            return resp.json()

    async def resolve_context(self, prompt: str) -> Dict[str, Any]:
        """Parse and resolve @, /, and # mentions."""
        async with httpx.AsyncClient(timeout=10.0) as client:
            resp = await client.post(
                f"{self.base_url}/api/v1/context",
                json={"prompt": prompt}
            )
            resp.raise_for_status()
            return resp.json()

    async def run(self, prompt: str, auto_enhance: bool = True, output_format: str = "json") -> Dict[str, Any]:
        """Execute a prompt with OpenCode agent and tools."""
        async with httpx.AsyncClient(timeout=600.0) as client:
            resp = await client.post(
                f"{self.base_url}/api/v1/run",
                json={
                    "prompt": prompt,
                    "auto_enhance": auto_enhance,
                    "format": output_format
                }
            )
            resp.raise_for_status()
            return resp.json()
