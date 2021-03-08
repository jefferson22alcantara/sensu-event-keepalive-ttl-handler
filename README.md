# Sensu Go Handler Event KeepAlive By TTL 
![Go Test](https://github.com/jefferson22alcantara/sensu-event-keepalive-ttl-handler/workflows/Go%20Test/badge.svg)
[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/jefferson22alcantara/sensu-event-keepalive-ttl-handler)

The Sensu Go Handler Event KeepAlive By TTL is a [Sensu Event Handler][2] which manages
Sensu Events  Generede by ad hoc Checks that still alerting after TTL expired as message log  **`"level":"warning","msg":"check TTL expired"`** and **`"output": "Last check execution was xX seconds ago"`** .The purposes of handler is [Sensu][1] remove that events after number of occurrences trigged.


## Installation

Download the latest version of the sensu-event-keepalive-ttl-handler from [releases][3],
or create an executable script from this source.

From the local path of the sensu-event-keepalive-ttl-handler repository:
```
go build -o /usr/local/bin/sensu-event-keepalive-ttl-handler main.go
```

## Configuration

Example Sensu Go handler definition:

```yml
type: Handler
api_version: core/v2
metadata:
  name: event-keepalive-ttl-handler
  namespace: default
spec:
  type: pipe
  command: sensu-event-keepalive-ttl-handler -k APYKEY -h http://localhost:8080 -m 10
  env_vars:
  - API_JEY=""
  - SENSU_MAX_OCCURRENCE=10
  - SENSU_URL=http://localhost:8080
  timeout: 10
  runtime_assets:
  - jefferson22alcantara/sensu-event-keepalive-ttl-handler
  filters:
  - is_incident
```

Example Sensu Go check definition:

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: dummy-app-healthz
  namespace: default
  annotations:
    sensu.io/plugins/sensu-event-keepalive-ttl-handler/conf/enable: "true"
spec:
  command: check-http -u http://localhost:8080/healthz
  subscriptions:
  - dummy
  handlers:
  - event-keepalive-ttl-handler
  interval: 60
  publish: true
```


## Usage Examples

Help:
```


The Sensu Go Remove events as keepalive

Usage:
  sensu-event-keepalive-ttl-handler [flags]
  sensu-event-keepalive-ttl-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -k, --apikey string        The Sensu Api Key , use default from API_JEY env var
  -h, --help                 help for sensu-event-keepalive-ttl-handler
  -m, --max_occurrence int   The Max event Occurence after ttl to remove event   , use default from SENSU_MAX_OCCURRENCE env var
  -u, --sensu_url string     The Sensu Api URL  , use default from SENSU_URL env var

Use "sensu-event-keepalive-ttl-handler [command] --help" for more information about a co

```

**Note:** Make sure to set the `"sensu.io/plugins/sensu-event-keepalive-ttl-handler/conf/enable": "true",` annotaition on checks. 

### Argument Annotations

All arguments for this handler are tunable on a per entity or check basis based on annotations.  The
annotations keyspace for this handler is `"sensu.io/plugins/sensu-event-keepalive-ttl-handler/conf/enable": "true",`. 

#### Examples

To change the team argument for a particular check, for that checks's metadata add the following:

```yml
type: CheckConfig
api_version: core/v2
metadata:
  annotations:
    sensu.io/plugins/sensu-event-keepalive-ttl-handler/conf/enable: "true"
[...]
```


### Asset creation

The easiest way to get this handler added to your Sensu environment, is to add it as an asset from Bonsai:

```sh
sensuctl asset add jefferson22alcantara/sensu-event-keepalive-ttl-handler
```

See `sensuctl asset --help` for details on how to specify version.

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/sensu/sensu-go
[2]: https://docs.sensu.io/sensu-go/5.0/reference/handlers/#how-do-sensu-handlers-work
[3]: https://github.com/jefferson22alcantara/sensu-event-keepalive-ttl-handler/releases
