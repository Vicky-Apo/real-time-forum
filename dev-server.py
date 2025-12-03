#!/usr/bin/env python3
"""
Simple development server for SPA with proper MIME types and routing
"""

import http.server
import socketserver
import os
from pathlib import Path

PORT = 3000
DIRECTORY = "client"

class SPAHandler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=DIRECTORY, **kwargs)

    def end_headers(self):
        # Add proper MIME types
        if self.path.endswith('.js'):
            self.send_header('Content-Type', 'application/javascript')
        elif self.path.endswith('.css'):
            self.send_header('Content-Type', 'text/css')
        elif self.path.endswith('.json'):
            self.send_header('Content-Type', 'application/json')

        # Disable caching for development
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
        self.send_header('Pragma', 'no-cache')
        self.send_header('Expires', '0')

        super().end_headers()

    def do_GET(self):
        # Remove leading slash
        path = self.path.lstrip('/')

        # If path is empty, serve index.html
        if not path or path == '/':
            path = 'index.html'

        # Build full file path
        full_path = os.path.join(DIRECTORY, path)

        # Check if file exists
        if os.path.isfile(full_path):
            # File exists, serve it normally
            super().do_GET()
        else:
            # File doesn't exist, serve index.html (SPA routing)
            self.path = '/index.html'
            super().do_GET()

if __name__ == "__main__":
    with socketserver.TCPServer(("", PORT), SPAHandler) as httpd:
        print(f"🚀 Development server running at http://localhost:{PORT}")
        print(f"📁 Serving from: ./{DIRECTORY}")
        print("📝 SPA routing enabled (all routes serve index.html)")
        print("Press Ctrl+C to stop\n")

        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\n\n👋 Server stopped")
