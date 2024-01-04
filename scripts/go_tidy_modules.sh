#!/bin/bash

# execute from root of project

cd engine/admin-api
go mod tidy

cd ../k8s-manager
go mod tidy

cd ../nats-manager
go mod tidy
