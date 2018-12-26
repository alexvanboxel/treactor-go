# Reactor

Reactor is a microservice designed to test and experiment with observability of microservices. You can play with it
on your machine in `local` more, but it gets interesting when you deploy it on a Kubernetes cluster or an Istio mesh.

In cluster mode you can have a bit over 120 microservices (atoms in the mendeleev table and a few more). You can change
the behaviour of the cloud of microservices by giving `reactor` some interesting `molecules`.

### Motivation

I created this microservice to let me inspect what happens inside the kubernetes cluster or a mesh. I've noticed in the
past that not all the scenario's worked that well together, certainly when you go a bit beyond the happy path. Well this
microservice will enable me in an easy way to *make reproducibles* that other people can run. I hope it's also helpful
for other people to learn the *learn* technologies like tracing, kubernetes networking, logging, istio, etc...

This reactor is for now focused on Google Cloud, as that's the cloud I live on, but once reactor is a bit more stable
I'm open for other people to add other providers.

Personally I will start using (and extending) reactor in bug reports, articles and talks. Let's see if this will get
as popular as httpbin is ;).

### Features

* Create cascading microservice calls for one origin call
* Configure tracing and logging
* Inspect headers, every level in the call hierarchy

### Example

`[H]*2[O]`

## Installation

### Pre-Requirement

Create a directory `tmp` and `work` in this repo. Don't worry, they are in the `.gitignore` so you do not accidentally
check them in.

Create a service account for *your* project (eg `my-research`). That will be used local as well on cluster. Download
the JSON private key and copy it to `work/my-research/service-account.json`. Read the documentation to see how to
create the service account:  [https://cloud.google.com/iam/docs/creating-managing-service-account-keys].

### Local

Build the `reactor`

`go install ./cmd/reactor/`

Set the environmental variables.

```
export PORT=3330
export REACTOR_NAME=reactor-api
export REACTOR_VERSION=1
export REACTOR_MODE=local
export GOOGLE_APPLICATION_CREDENTIALS=/Users/me/path/work/my-research/service-account.json
export GOOGLE_PROJECT_ID=my-research
```

Execute.

`reactor`

Test

http://localhost:3330/rr/split?molecule=[H]^2[O]

Go to the Cloud Console, select *Trace*.

### Cluster

*Not yet tested/supported*
`gcloud builds submit --config=cb.yaml .`


## Specification

### Environment Variables

NAME | Description | Default
---- | ----------- | -------
PORT | Port |
APP_NAME | Application name |
APP_VERSION | Application version |
REACTOR_MODE | Reactor mode (local, k8s) | local
REACTOR_DEBUG | Debug Mode (0, 1) | 1
REACTOR_PROFILE | Stackdriver Profiler (0, 1) | 0
REACTOR_TRACE_PROPAGATION | OpenCensus propagator (stackdriver, b3)  | stackdriver
REACTOR_TRACE_INTERNAL | Reactor adds internal spans (0, 1) | 1



### Molecule spec


```
S    [H,x=1,y=2,z=3]*2[O,x=1,y=2,z=3],x=1,y=2,z=3
A     H,x=1,y=2,z=3
A                      O,x=1,y=2,z=3
```


```
S    2[5[Ur,log:1,xyz:4]^5[C,log:1,xyz:4]],x:1,y:2
O1     5[Ur,log:1,xyz:4]^5[C,log:1,xyz:4]
A        Ur,log:1,xyz:4
A                          C,log:1,xyz:4
```

## Test Scenarios

https://github.com/wg/wrk

`wrk -t12 -c400 -d30s "http://<yourip>/rr/split?molecule=[H]^2[O]"`

