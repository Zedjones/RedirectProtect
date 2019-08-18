#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"

pushd routes > /dev/null
COMPOSEPATH=$SCRIPTPATH go test -cover
popd > /dev/null
