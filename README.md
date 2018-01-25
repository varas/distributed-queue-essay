# Essay on a distributed in-memory queue

Quick essay on building a **replicated** in-memory queue.

> Just an ESSAY, not a working solution.

This is just a quick micro-service skeleton.

No use-case in mind, so I won't add distributed consistency and coordination management.

Depending on the given use-case we may go for a CRDT based solution to avoid coordination.

Nevertheless in-memory *queue* is modeled to handle operation replication.

## Approach

A node is run on 2 TCP ports:

* `queue operations` port
* `cluster operations` port

*Cluster operations* are interfaced under `cluster/` folder.

*Queue operations* are almost implemented on `queue/` folder.

I used plain TCP for simplicity, but HTTPS, gRPC, whatever could be used for both communications.

> Also *Go-Kit* could be used as a general good practice micro-service design.

### Cluster approach

I avoided cluster management is bounded to the `cluster` package.

`Cluster` models cluster membership management via:

* 2 request methods
* 2 receiver channels.

`Publisher` is supposed to manage synchronization in background mode.
It also allows domain code to publish queue changes to the cluster.

### Node approach

*Node daemon logic* is under `node/`.

Several connection handlers are used to concurrent update the queue.

Queue is managed via command-handler pattern to easily support replication operations.

## Run

`make run`

> This essay is not aimed to be compiled.

## Trace queue requests

To-do: add tracing capabilities `go get github.com/jaegertracing/jaeger`

## Profile node

To profile CPU usage: `make profile-cpu`

Memory profile could also be added
