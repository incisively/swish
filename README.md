# Swish

A simple tool to help with zero-downtime Docker deployments. It a switch, ya bish.


## Because

If you assign a public port to a Docker container, you cannot do zero-downtime deployment because you have to kill the old container to free the port, and then wait for the new one to boot up before service resumes.

You could use `iptables` but fuck that.

## Solution

Swish listens to any number of public ports, and reverse proxies each to any target host and port. Target can be updated instantly using a simple HTTP request.

Swish works well manually or as part of a deployment workflow.

## Usage

When you start Swish, the admin API will come up on `0.0.0.0:8999` (you can change this using the `-bind` option). It has two methods:

#### Summary

Any `GET` request will return a plaintext summary of current listeners and targets.

#### Update

Any `POST` request with `listen` and `target` parameters will tell Swish to:

- Start listening on `listen` if it isn't already
- Forward any requests for `listen` to `target`

`listen` and `target` must be in the format `host:ip`, e.g. `localhost:12345`. To listen on all interfaces, you can drop the host portion of `listen` (i.e. `:9001` implies `0.0.0.0:9001`).

## Example

#### Background

Suppose we want to expose some services on fixed ports in the `9000-9100` range.

First off, start Swish, assigning the entire public port range you'd like to use to a Swish container, plus the one for Swish's API (that's `8999` by default):

```sh
$ docker run -d -p 8999-9100:8999-9100 incisively/swish
```

We only have to do this once, since one Swish instance will proxy multiple ports.

> You could just run Swish directly on the host, but a Docker container is easier than writing Upstart scripts or daemonizing a Go program.

#### Expose a service

Let's bring up our first service, which happens to expose port `8080`. We use `-P` to let Docker assign an ephemeral host port:

```sh
docker run -dP --name awesome-service awesome-service:1.0.0
```

Now let's make that container's port `8080` available to the outside world on the fixed port `9000`:

```sh
# n.b. you must set $DOCKER_GATEWAY_IP yourself.
IP=$DOCKER_GATEWAY_IP
PORT=$(docker port awesome-service 8080 | cut -d ':' -f 2)

# update Swish
curl localhost:8999 -d "listen=:9000&target=$IP:$PORT"
```

Done. Sweet deal.

#### Replace the service

Let's say later on we replace our `awesome-service` container with a new version:

```sh
$ docker rename awesome-service awesome-service-old
$ docker run -dP --name awesome-service awesome-service:1.0.1
```

Docker will assign a different host port to the new container, but all we need to do is run the same code as before to tell Swish to switch to the new container:

```sh
# n.b. you must set $DOCKER_GATEWAY_IP yourself.
IP=$DOCKER_GATEWAY_IP
PORT=$(docker port awesome-service 8080 | cut -d ':' -f 2)

# update Swish
curl localhost:8999 -d "listen=:9000&target=$IP:$PORT"
```

All traffic to `:9000` now goes to the new container, and we didn't drop any requests. We're free to kill and remove the old one:

```sh
$ docker kill awesome-service-old
$ docker rm awesome-service-old
```

# TODO

- Ability to delete listeners
- Better testing
