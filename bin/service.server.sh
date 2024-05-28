#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../app && pwd )

export JELLYFIN_kestrel__socket=true
export JELLYFIN_kestrel__socketPath=/var/snap/jellyfin/current/socket
export JELLYFIN_kestrel__socketPermissions=0777

LIBS=$(echo ${DIR}/usr/lib/*-linux-gnu*)
LIBS="$LIBS:$(echo ${DIR}/lib/*-linux-gnu*)"
LIBS="$LIBS:$(echo ${DIR}/usr/lib)"
exec ${DIR}/lib/*-linux*/ld-*.so --library-path $LIBS $DIR/jellyfin/jellyfin --configdir /var/snap/jellyfin/current/config --datadir /var/snap/jellyfin/current/data --cachedir /var/snap/jellyfin/current/cache --ffmpeg /snap/jellyfin/current/bin/ffmpeg.sh
