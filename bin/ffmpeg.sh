#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../app && pwd )
LIBS=$(echo ${DIR}/usr/lib/*-linux-gnu*)
LIBS="$LIBS:$(echo ${DIR}/usr/lib/jellyfin-ffmpeg)/lib"
exec ${DIR}/lib/*-linux*/ld-*.so --library-path $LIBS ${DIR}/usr/lib/jellyfin-ffmpeg/ffmpeg "$@"
