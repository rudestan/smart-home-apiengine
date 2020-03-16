package alexakit

import (
    "os"
    "smh-apiengine/pkg/amqp"

    "github.com/spf13/cast"
)

const (
    RmqHost       = "localhost"
    RmqPort       = 5672
    RmqLogin      = "guest"
    RmqPassword   = "guest"
    RmqExchange   = "alexa_sync"
    RmqQueue      = "alexa.responses"
    RmqRoutingKey = "alexa.response.json"
)

const (
    EnvRmqHost = "SMH_PROXY_RMQ_HOST"
    EnvRmqPort = "SMH_PROXY_RMQ_PORT"
    EnvRmqLogin = "SMH_PROXY_RMQ_LOGIN"
    EnvRmqPassword = "SMH_PROXY_RMQ_PASSWORD"
    EnvRmqExchange = "SMH_PROXY_RMQ_EXCHANGE"
    EnvRmqQueue = "SMH_PROXY_RMQ_QUEUE"
    EnvRmqRoutingKey = "SMH_PROXY_RMQ_ROUTING_KEY"
)

func NewConfigFromEnv() *amqp.Config {
    config := amqp.Config{
        Host:       getEnvVar(EnvRmqHost, RmqHost),
        Port:       cast.ToInt(getEnvVar(EnvRmqPort, cast.ToString(RmqPort))),
        Login:      getEnvVar(EnvRmqLogin, RmqLogin),
        Password:   getEnvVar(EnvRmqPassword, RmqPassword),
        Exchange:   getEnvVar(EnvRmqExchange, RmqExchange),
        Queue:      getEnvVar(EnvRmqQueue, RmqQueue),
        RoutingKey: getEnvVar(EnvRmqRoutingKey, RmqRoutingKey),
    }

    return &config
}

func getEnvVar(name string, defaultVal string) string  {
    if value, exists := os.LookupEnv(name); exists {
        return value
    }

    return defaultVal
}