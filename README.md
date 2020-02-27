### Alexa Broadlink Devices Control

Application for controlling Broadlink devices using incoming Requests from preconfigured Alexa skill. Instead of HTTP
server that handles Alexa API POST requests this application is using RabbitMQ for consuming the pre-stored json 
requests from the queue. 

This application can be started for example on RaspberryPi that is on the same network with Broadlink devices.

Since it is only part of the smart home implementation, there are few things still missing:

- Middleware application that can handle POST requests and publish them to the RabbitMQ queue, so the messages can be consumed by this application.
Raw POST request is coming in JSON format so it can be just pushed directly to the queue.
- One of the two options: 
    - Configured HTTPS server with either valid certificate or selfsigned certificate that can handle incoming requests
    from Alexa API using the Middleware application from above
    - Middleware application which can be uploaded and executed in Amazon Lambda

There is already Middleware Go application that uses the first option from above and can be hosted on some VPS. It will
be published on github as well.  