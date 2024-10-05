#!/bin/bash

output=$1
if [[ -z "$output" ]]; then
  echo "usage: $0 <output-name>"
  exit 1
fi

packages=$2
if [[ -z "$packages" ]]; then
  echo "usage: $1 <package-name>"
  exit 1
fi

platforms=("darwin/amd64" "linux/amd64" "windows/amd64") #"darwin/386" "dragonfly/amd64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/ppc64" "linux/ppc64le" "linux/mips" "linux/mipsle" "linux/mips64" "linux/mips64le" "netbsd/386" "netbsd/amd64" "netbsd/arm" "openbsd/386" "openbsd/amd64" "openbsd/arm" "solaris/amd64" "windows/386" "windows/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    name=$output

    if [ $GOOS = "windows" ]; then
        name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o ./build/$GOOS/$GOARCH/$name $packages
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'$GOOS/$GOARCH
        exit 1
    fi
    echo $GOOS/$GOARCH complete
done

