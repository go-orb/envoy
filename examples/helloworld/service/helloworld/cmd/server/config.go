package main

//nolint:gochecknoglobals
var config = `{
    "client": {
        "middlewares": [
            {
                "name": "retry"
            },
            {
                "name": "log"
            }
        ]
    },
    "registry": {
        "plugin": "kvstore",
        "kvstore": {
            "plugin": "natsjs",
            "servers": [
                "nats://nats:4222"
            ]
        }
    }
}`
