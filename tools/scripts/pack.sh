#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

# cd To root of project
cd $DIR/../../
 
mkdir -vp packed/configs
cp -vfR out/mcli packed
cp -vfR configs/mcli.json packed/configs
cp -vfR configs/packages_conf.json packed/configs
cp -vfR configs/processes_conf.json packed/configs


