# K8s manager

- [K8s manager](#k8s-manager)
  - [Description](#description)
  - [gRPC](#grpc)
  - [Kubernetes](#kubernetes)

### Description

This service is part of the Engine, exposes a gRPC service to encapsulate all Kubernetes related features. The only
service that is going to call this gRPC is the Admin API service when need to create new Kubernetes resources.

### gRPC

The Protobuf file and the code generated are within `proto` folder.

To generate the code from the `.proto` file run the following command.

```bash
./scripts/generate_proto.sh
```

We expose the following service in the gRPC server:

- **VersionService**: Is intended to control the versions lifecycle with the following functions.
  - Start
  - Stop
  - Publish
  - Unpublish
  - UpdateConfig

### Kubernetes

This server uses the official kubernetes sdk `client-go` to interact with the cluster.
