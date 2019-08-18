#!/bin/bash

pushd routes > /dev/null
go test
popd > /dev/null
