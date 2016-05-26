#!/usr/bin/env bash

set -e

source ././../../docker/Constant.rc

# clean dangling images
# mac not support
# docker images -q --filter "dangling=true"|awk '{print $0}' | xargs -t -i docker rmi -f {}
docker images -q --filter "dangling=true"| xargs docker rmi -f

#package go program return dir now
Cluster="gopush"
TarFile=${Cluster}-${Version}.tar.gz
echo "Dir now:" `pwd`

#compile go program
echo -e "Will build go program use docker container from image: "$ImageBuild"..."
docker run -v ${VolumePath}:${VolumePath} $ImageBuild bash -c "go build -a -ldflags '-s' gopush \
    && mv gopush ${ProjectPath}/release/tar/pack-${Cluster}/transfer/bin"

# package go program return dir now
echo -e "Will package go program into tar file...\nDir now:" `pwd`

mkdir -p output/${Cluster}-${Version}
cp -a transfer/ output/${Cluster}-${Version}

if [ -f output/$TarFile ]; then
    rm output/$TarFile
fi

# 需要cd到output目录.
cd output
tar -czf $TarFile ${Cluster}-${Version}/
cd -

# upload file to hosts.
cp output/$TarFile $FileServer
echo "URL: http://docker.alishui.com/tmp/${Cluster}-${Version}.tar.gz"



