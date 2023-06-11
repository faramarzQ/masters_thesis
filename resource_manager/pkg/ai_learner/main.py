from src.http_server.server import runServer

if __name__ == "__main__":
    runServer()

    # Python 3 server example
from http.server import BaseHTTPRequestHandler, HTTPServer
import time

hostName = "localhost"
serverPort = 8080