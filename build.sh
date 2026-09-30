#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

VERSION=$1
LDAP_VERSION=$2

if [[ -z "$LDAP_VERSION" ]]; then
    echo "usage $0 version ldap_version"
    exit 1
fi

BUILD_DIR=${DIR}/build/snap/app
${DIR}/apt.sh unzip wget
mkdir -p ${BUILD_DIR}
cd ${BUILD_DIR}
cp -r /usr ${BUILD_DIR}
cp -r /lib ${BUILD_DIR}
cp -r /bin ${BUILD_DIR}
cp -r /jellyfin ${BUILD_DIR}

ARCH_DIR=$(dirname usr/lib/*/ld*.so.*)
ln -s /snap/jellyfin/current/app/jellyfin/jellyfin.dll $ARCH_DIR/jellyfin.dll
ls -la $ARCH_DIR/jellyfin.dll

${DIR}/download.sh \
  https://repo.jellyfin.org/files/plugin/ldap-authentication/ldap-authentication_${LDAP_VERSION}.zip \
  ${DIR}/build/ldap-authentication.zip

mkdir -p $BUILD_DIR/plugins/LDAP-Auth
unzip ${DIR}/build/ldap-authentication.zip -d $BUILD_DIR/plugins/LDAP-Auth
ls -la $BUILD_DIR/plugins/LDAP-Auth

PLUGIN_ABI=$(sed -n 's/.*"targetAbi": *"\([0-9]*\)\..*/\1/p' $BUILD_DIR/plugins/LDAP-Auth/meta.json)
SERVER_MAJOR=${VERSION%%.*}
if [[ "$PLUGIN_ABI" != "$SERVER_MAJOR" ]]; then
    echo "ldap plugin ${LDAP_VERSION} targets jellyfin ${PLUGIN_ABI}.x, server is ${VERSION}"
    exit 1
fi
