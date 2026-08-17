#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

if [[ -z "$2" ]]; then
    echo "usage $0 app version"
    exit 1
fi

NAME=$1
VERSION=$2
ARCH=$(dpkg --print-architecture)

SNAP_DIR=${DIR}/build/snap

apt update
apt -y install squashfs-tools wget

cp -r ${DIR}/bin ${SNAP_DIR}
cp -r ${DIR}/config ${SNAP_DIR}
cp -r ${DIR}/meta ${SNAP_DIR}

echo "version: $VERSION" >> ${SNAP_DIR}/meta/snap.yaml
echo "architectures:" >> ${SNAP_DIR}/meta/snap.yaml
echo "- ${ARCH}" >> ${SNAP_DIR}/meta/snap.yaml
echo $VERSION > ${SNAP_DIR}/version

for f in \
    meta/snap.yaml \
    meta/gui/icon.png \
    meta/hooks/install \
    meta/hooks/configure \
    meta/hooks/pre-refresh \
    meta/hooks/post-refresh \
    bin/cli \
    bin/service.server.sh \
    bin/service.nginx.sh \
    bin/ffmpeg.sh \
    bin/ffprobe.sh \
    config/nginx.conf \
    config/jellyfin/config/network.xml \
    nginx/bin/nginx.sh \
    app/jellyfin/jellyfin \
    version; do
    if [[ ! -e ${SNAP_DIR}/${f} ]]; then
        echo "missing from snap: ${f}"
        exit 1
    fi
done

du -d10 -h $SNAP_DIR | sort -h | tail -100

PACKAGE=${NAME}_${VERSION}_${ARCH}.snap
echo ${PACKAGE} > ${DIR}/package.name
mksquashfs ${SNAP_DIR} ${DIR}/${PACKAGE} -noappend -comp xz -no-xattrs -all-root
mkdir ${DIR}/artifact
cp ${DIR}/${PACKAGE} ${DIR}/artifact

