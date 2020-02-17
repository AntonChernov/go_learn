#!/bin/bash
#export GOPATH=$HOME/go
#export PATH=$GOPATH/bin:$PATH


if test -z "$1" 
then
      echo "\$1 is empty"
else
      echo "Dir is $1"
      #export GOPATH=/Users/anton_chernov2/$1
      export GOPATH=$1
      export PATH=$PATH:$(go env GOPATH)/bin
      cd $GOPATH/
fi
