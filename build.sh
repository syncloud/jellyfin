#!/bin/sh -ex

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}

LDAP_VERSION=$1

BUILD_DIR=${DIR}/build/snap/app
apt update
apt install -y unzip
mkdir -p ${BUILD_DIR}
echo $VERSION > $BUILD_DIR/app.version
cd ${BUILD_DIR}
cp -r /usr ${BUILD_DIR}
cp -r /lib ${BUILD_DIR}
cp -r /bin ${BUILD_DIR}
cp -r /jellyfin ${BUILD_DIR}

ARCH_DIR=$(dirname usr/lib/*/ld*.so.*)
ln -s /snap/jellyfin/current/app/jellyfin/jellyfin.dll $ARCH_DIR/jellyfin.dll
ls -la $ARCH_DIR/jellyfin.dll

curl https://repo.jellyfin.org/releases/plugin/ldap-authentication/ldap-authentication_${LDAP_VERSION}.zip -o ${DIR}/build/ldap-authentication.zip
mkdir -p $BUILD_DIR/plugins/LDAP-Auth

unzip ${DIR}/build/ldap-authentication.zip -d $BUILD_DIR/plugins/LDAP-Auth

ls -la $BUILD_DIR/plugins/LDAP-Auth
