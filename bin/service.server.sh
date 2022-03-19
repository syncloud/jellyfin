#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

$DIR/jellyfin/jellyfin --datadir /var/snap/jellyfin/current/data --cachedir /var/snap/jellyfin/current/cache --ffmpeg /snap/jellyfin/current/app/usr/lib/jellyfin-ffmpeg/ffmpeg
fmpeg
