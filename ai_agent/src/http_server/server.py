from http.server import BaseHTTPRequestHandler, HTTPServer
from src.configs.configs import *
from src.services.rl_agent import *
from urllib.parse import parse_qs
import json 
from datetime import date
import logging


class requestHandler(BaseHTTPRequestHandler):
    """
        Request handler class which includes http API handler methods
    """
    def readBody(self):
        length = int(self.headers.get('content-length'))
        postBodyBytes = self.rfile.read(length)
        body = json.loads(str(postBodyBytes,"UTF-8"))
        return body

    def setSuccessHeaders(self):
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()

    def setInternalErrorHeaders(self):
        self.send_response(500)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()

    def do_POST(self):
        body = self.readBody()

        logging.basicConfig(level=logging.DEBUG, filename="src/storage/logs.log", filemode="a+",
            format="%(asctime)-15s %(levelname)-8s %(message)s")

        logging.info("------------------------------------------")
        logging.info("Run at step: %s", body["Step"])
        print("Run at step: ", body["Step"])
        logging.info("------------------------------------------")

        logging.info("Request: %s", body)

        response = {}
        try:
            response = runReinforcementLearning(body)
        except Exception as e:
            logging.error("error", exc_info=True)
            print("Error at step: ", body["Step"])
            self.setInternalErrorHeaders()
            return

        responseString = json.dumps(response)
        logging.info("Response: %s", responseString)

        self.setSuccessHeaders()
        self.wfile.write(responseString.encode(encoding='utf_8'))

def runServer():
    """
    Runs an http server exposing APIs

    """
    print("Server started http://%s:%s" % (hostName, serverPort))
    webServer = HTTPServer((hostName, serverPort), requestHandler).serve_forever()

    try:
        webServer.serve_forever()
    except KeyboardInterrupt:
        pass

    webServer.server_close()
    print("Server stopped.") 