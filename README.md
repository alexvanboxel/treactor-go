> Reactor is using a patched version of google-cloud-go. The patch is submitted but is waiting a release, do not
> update the dependencies (dep). This block will be removed when it's safe to update.

# Reactor

Reactor is a microservice designed to test and experiment with observability of microservices. You can play with it
on your own machine in `local` more, but it gets interesting when you deploy it on a Kubernetes cluster or an Istio mesh.

In cluster mode you can have a bit over 120 microservices (atoms in the mendeleev table). You control `reactor` some
giving it some interesting `molecules`.

### Motivation

I created this microservice to let me inspect what happens inside the kubernetes cluster or a mesh. I've noticed that
not all the scenario's and products worked that well together, certainly when you go a bit beyond the happy path. This
microservice will enable me to *make reproducibles*. I hope it's also helpful for other people to *learn*  technologies
like tracing, kubernetes networking, logging, istio, etc...

Reactor is for now focused on Google Cloud, as that's the cloud I live on, but once reactor is a bit more stable I'm open
for other people to add other providers.

I will start using (and extending) reactor for bug reports, articles and talks. Let's see if this will get
as popular as httpbin is ;).

### Features

* Create cascading microservice calls from one origin call
* Configure tracing and logging
* Inspect headers, every level in the call hierarchy

### Example

To see what reactor does lets take a very simple example. Reactor works by splitting molecules:

`[[H]]^2[O]`

Everything between the brackets `[]` will be a call to the next microservice. The brackets can be prefixed with a number
and an optional parameter (`s` or `p`), this tells reactor how many times the service needs to be called and how (
sequential or parallel). Multiple calls can be make by appending them using `^` (sequential) or `*` (parallel).

Depending the content of the bracket the call will be different. If reactor detects an atom a call to the corresponding
atom service will be made. But if reactor detects another sub-molecule it calls the next orbit and apply the same
logic till only atoms are left. So the example above will result in:

`http://reactor-api/rr/split?molecule=[[H]]^2[O]`

calling

* `http://orbit-1/rr/split?molecule=[H]`
* `http://atom-o/rr/atom?atom=O`
* `http://atom-o/rr/atom?atom=O`

*orbit* will split the molecule [H] (ok, this looks strange, but each bracket is a layer) into it's atoms, in this
case only 1 `H`:

* `http://atom-h/rr/atom?atom=H`

Try the local installation, to see how it looks in the trace (this will make it more clear).

## Installation

### Pre-Requirement

Create a directory `tmp` and `work` in this repo. Don't worry, they are in the `.gitignore` so you do not accidentally
check them in.

Create a Google Cloud project, enable:

* Tracing
* Logging
* Profile
* Debugging
* Cloud Build
* Container Registry

> Replace project-name with your project name in the rest of the instructions

Create a service account for *your* project (eg `project-name`) with the appropriate roles (see above) .
That will be used local as well on cluster. Download the JSON private key and copy it to
`work/project-name/service-account.json`. Read the documentation to see how to
create the service account:  [https://cloud.google.com/iam/docs/creating-managing-service-account-keys].

### Local

Build the `reactor` from source

`go install ./cmd/reactor/`

Set the environmental variables. Replace `project-name` project with your own

```
export PORT=3330
export REACTOR_NAME=reactor-api
export REACTOR_VERSION=1
export REACTOR_DEBUG=1
export REACTOR_PROFILE=0
export REACTOR_MODE=local
export GOOGLE_APPLICATION_CREDENTIALS=/Users/me/path/work/project-name/service-account.json
export GOOGLE_PROJECT_ID=project-name
```

Fire up the reactor

`reactor`

And test it by calling:

http://localhost:3330/rr/split?molecule=[[H]]^2[O]

Go to the Cloud Console, select *Trace*.

### Kubernetes

*Not yet fully tested/supported*

`gcloud builds submit --config=cb.yaml .`

`go install ./cmd/rrprep/`

`kubectl label namespace default istio-injection=disabled --overwrite`

### Istio

*Not yet fully tested/supported*

`kubectl label namespace default istio-injection=enabled --overwrite`

## Specification

### Reactor Prepare (rrprep)

Prepares the Kubernetes files, from the templates

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

