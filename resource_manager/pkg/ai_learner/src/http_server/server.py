from http.server import BaseHTTPRequestHandler, HTTPServer
from src.configs.configs import *
from src.services.rl_agent import *
from urllib.parse import parse_qs
import json 


class requestHandler(BaseHTTPRequestHandler):
    def setHeaders(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

    def readBody(self):
        length = int(self.headers.get('content-length'))
        postBodyBytes = self.rfile.read(length)
        body = json.loads(str(postBodyBytes,"UTF-8"))

        return body

    def do_POST(self):
        print("Post request")

        body = self.readBody()
        print(body)

        self.setHeaders()

def runServer():
    """
    Runs an http server exposing APIs

    """
    webServer = HTTPServer((hostName, serverPort), requestHandler).serve_forever()
    print("Server started http://%s:%s" % (hostName, serverPort))

    try:
        webServer.serve_forever()
    except KeyboardInterrupt:
        pass

    webServer.server_close()
    print("Server stopped.")