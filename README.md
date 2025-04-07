# envoy

This repo contains the source code for the Envoy plugin for the [go-orb](https://github.com/go-orb/go-orb) framework.

## Demo usage

```sh
git clone https://github.com/go-orb/envoy.git
cd envoy/examples/helloworld
docker compose up -d --build
```

### Bench envoy

```sh
docker compose run --rm fortio load -payload '{}' -H "Content-Type: application/json" -qps -1 -gomaxprocs 12 -c 24 -t 15s http://envoy2:10000/hello.v1.Hello/Hello
```

### Bench the plugin

```sh
docker compose run --rm fortio load -H "Content-Type: application/json" -qps -1 -gomaxprocs 12 -c 24 -t 15s http://envoy:10000/hello.v1.Hello/Hello
```

### Bench traefik

```sh
docker compose run --rm fortio load -payload '{}' -H "Content-Type: application/json" -qps -1 -gomaxprocs 12 -c 24 -t 15s http://traefik:10000/hello.v1.Hello/Hello
```

### Bench backends

```sh
docker compose run --rm fortio load -payload '{}' -H "Content-Type: application/json" -qps -1 -gomaxprocs 12 -c 24 -t 15s http://helloworld:10000/hello.v1.Hello/Hello
```

## License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

## Copyright

Copyright 2025 René Jochum
