
![tardigrade-connector-framework-arch](https://raw.githubusercontent.com/utropicmedia/tardigrade-connector-framework/master/images/conector-framework-arch.png)


# Tardigrade Connector Framework

The Tardigrade Connector Framework[TCF] is a basic set of utility methods and operations to provide a consisten approach to integrating and orchestrating interactions between data sources, endpoints, and the Tardigrade network.  Some of the aspects we address are:
  
  - Buffering/Resource management
  - Abstraction
  - Data Transformation
  - Configuration
  - Authentication
  - Operations
  - Logging


### Installation

The SCF requires [Go](https://golang.org/) to run.

Install the dependencies and devDependencies and start the server.

Install storj-uplink go package, by running:


``` sh
$ go get storj.io/storj/lib/uplink
```

## TODO

  - Add Logging config
  - Add Orchestration config
  - Add buffer tuning parameters
