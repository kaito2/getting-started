#! /bin/bash -e

if [ $# != 1 ]
then
  echo "Usage: $0 dir"
  exit 1
fi

if [ -e $1 ]; then
    echo "Directory $1 already exists."
fi

mkdir -p $1
cp templates/README.md $1
