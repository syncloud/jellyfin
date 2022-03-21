#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
LIBS=$(echo ${DIR}/usr/lib/*-linux-gnu*)
${DIR}/lib/*-linux*/ld-*.so --library-path $LIBS ${DIR}/usr/lib/jellyfin-ffmpeg/ffmpeg "$@"
"$@"
