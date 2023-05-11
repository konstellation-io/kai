#!/bin/bash

cd adapter/dataloader
go run github.com/vektah/dataloaden UserLoader string *github.com/konstellation-io/kai/engine/admin-api/domain/entity.User
