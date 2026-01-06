import http.server
import socketserver
import os

PORT = 8080
DIRECTORY = "web/public"

class SquavaHandler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=DIRECTORY, **kwargs)

    def end_headers(self):
        # We'll remove COOP/COEP for now as they trigger strict blocking in Firefox
        # for Workers and imported scripts unless everything is perfectly configured.
        # Standard Go WASM does not require SharedArrayBuffer/Cross-Origin Isolation.
        
        # Same-origin resource policy is generally safe and good practice.
        self.send_header("Cross-Origin-Resource-Policy", "same-origin")
        
        # Ensure .gz files are handled correctly by the browser.
        if self.path.endswith(".gz"):
            self.send_header("Content-Encoding", "gzip")
        
        # Add a header to prevent caching during development.
        self.send_header("Cache-Control", "no-cache, no-store, must-revalidate")
        
        super().end_headers()

    def guess_type(self, path):
        # Explicitly handle WASM and JS, including gzipped versions.
        if path.endswith(".wasm") or path.endswith(".wasm.gz"):
            return "application/wasm"
        if path.endswith(".js") or path.endswith(".js.gz"):
            return "application/javascript"
        
        # Fallback to the default guesser for other types (html, css, etc.)
        base_guess = super().guess_type(path)
        
        # Ensure we don't return None for common types if the OS mime-db is missing them.
        if not base_guess:
            if path.endswith(".html"): return "text/html"
            if path.endswith(".css"): return "text/css"
            
        return base_guess

if __name__ == "__main__":
    script_dir = os.path.dirname(os.path.abspath(__file__))
    os.chdir(script_dir)
    
    socketserver.TCPServer.allow_reuse_address = True
    with socketserver.TCPServer(("", PORT), SquavaHandler) as httpd:
        print(f"Serving Squava at http://localhost:{PORT}")
        print(f"Root directory: {os.path.join(script_dir, DIRECTORY)}")
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\nShutting down server.")
