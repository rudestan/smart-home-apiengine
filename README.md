### Broadlink Devices Api Engine

Broadlink devices Api engine controls Broadlink Wi-Fi smart home Devices using incoming Requests either in JSON format 
of Amazon Alexa Api or directly by command or scenario id. The application works as a web server and serves requests 
view HTTP or HTTPS protocol. Additionally the commands or scenarios can be also executed via command line.

This application can be started for example on RaspberryPi that is on the same network with Broadlink devices.

Web server has the following end points:

1. ``GET`` ``/run/command/{commandId}`` - where ``{commandId}`` is the id of the command from the configuration JSON file
2. ``GET`` ``/run/scenario/{scenarioId}`` - where ``{scenarioId}`` is the id of scenario from the configuration JSON file
3. ``POST`` ``/run/intent`` - The ``POST`` data is the raw Json data coming from Amazon Alexa API, so the web server
can handle Alexa requests directly (HTTPS mode with valid certificate and key required by Amazon Alexa API) using
this endpoint, or this request can be proxified using RabbitMQ by rmqproxy, alexalistener tools from this project.

#### Usecases

##### Standalone HTTP

This application runs built in web server and to run it the only thing you need is configuration file with needed
devices and commands. 

To start serving requests and control the devices just run the application as following:

``apiengine -c ./path/to/your/config.json``

By default web server will be started on ``127.0.0.1:8787`` listening to HTTP requests

##### Standalone HTTPS

The application also supports HTTPS which is required by Amazon Alexa Api for example, so let's say our IP is 
publicly available, or some domain configured to point to IP where the app is planned to run. In this case
in the Amazon Alexa Skill configuration we can just set the skill endpoint to ``htts://<yourdomain|ip>/run/intent``
and run server as following:

``apiengine -c ./path/to/your/config.json serve --proto https --address "<yourdomain|ip>" --tls-cert "<pathtocert>" --tls-key "<pathtokey>"``

According to Amazon Alexa documentation, the certificate can be self signed if the skill is not planned for production use.

##### Whole tool set

In case you have your own public web server with https and you do not want to expose the static IP of some local machine in your
home network RabbitMQ can be used as a bridge between Amazon Alexa API and this Broadlink Api engine.
However you need some publicly available RabbitMQ server for working with messages.
Let's say you have one.

The infrastructure in this case can be the following:

1. On your publicly available server with HTTPS you should have some endpoint which will just push the incoming requests from Alexa Api
to some queue. Or alternatively ``requestproxy`` from the toolset can be used. It can run under ```nginx``` which 
configured as reverse proxy for that requests.

2. In your local network have some machine organized as server that can run ``rmqproxy`` from this tool set to consume to the
incoming messages from the queue above. Good case to use your RaspberryPi running under raspbian for example.

3. On the same machine (can be also separate machine) in your local network run this app in server mode to be able to serve the
requests that coming from ``rmqproxy``.

With this setup you should be able to serve the requests from Alexa Api and control your broadlink devices with voice. 
The latency is minimal while using RabbitMQ and the tools described above. Additionally you can run the ```webui``` from toolset and also control 
the devices via some touch interface (e.g. tablet)


ToDo:

- add some new apis for reloading the config
- add tests
- add some CI
- Write documentation
- visualise the toolset infrastructure