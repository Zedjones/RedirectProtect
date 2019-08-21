#!/bin/bash

pushd routes > /dev/null
go test -cover
popd > /dev/null

pushd internal > /dev/null
go test -cover
popd > /dev/null
