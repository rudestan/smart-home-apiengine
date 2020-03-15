### Smart Home RabbitMQ Proxy App 

Simple application that listens configured RabbitMQ server and forwards the payload to some preconfigured web server. 

Developed for usage in conjunction with ``apiengine`` from smart home toolset for forwarding Amazon Alexa Api requests, but can be easily 
used anywhere if there is need to forward rmq messages from the queue to the web endpoint.

#### Cli options

- ``--host``, ``-t``- RabbitMQ Host (default: "localhost")
- ``--port``, ``-p`` - RabbitMQ Host (default: 5672)
- ``--login``, ``-l`` - RabbitMQ Login (default: "guest")
- ``--password``, ``-s`` - RabbitMQ Password (default: "guest")
- ``--exchange``, ``-e`` - RabbitMQ Exchange name (default: "alexa_sync")
- ``--queue``, ``-q`` - RabbitMQ Queue name (default: "alexa.responses")
- ``--rkey``, ``-r`` - RabbitMQ Queue name (default: "alexa.response.json")
- ``--endpoint``, ``-u`` - Endpoint where to post the payload using POST method (default: "http://localhost:8787/run/intent")
- ``--log`` - Log file for logs output
- ``--help``, ``-h`` - show help (default: false)
