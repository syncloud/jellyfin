#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

export JELLYFIN_kestrel__socket=true
export JELLYFIN_kestrel__socketPath=/var/snap/jellyfin/current/socket

$DIR/jellyfin/jellyfin --datadir /var/snap/jellyfin/current/data --cachedir /var/snap/jellyfin/current/cache --ffmpeg /snap/jellyfin/current/app/usr/lib/jellyfin-ffmpeg/ffmpeg
fmpeg
