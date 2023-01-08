#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

ARCH=$(uname -m)
DOWNLOAD_URL=https://github.com/syncloud/3rdparty/releases/download/
LDAP_VERSION=16.0.0.0

BUILD_DIR=${DIR}/build/snap
mkdir -p $BUILD_DIR

apt update
apt -y install wget unzip

cd ${DIR}/build
wget -c --progress=dot:giga ${DOWNLOAD_URL}/nginx/nginx-${ARCH}.tar.gz
tar xf nginx-${ARCH}.tar.gz
mv nginx ${BUILD_DIR}

#binary
#wget -c --progress=dot:giga https://repo.jellyfin.org/releases/plugin/ldap-authentication/ldap-authentication_${LDAP_VERSION}.zip -O ldap-authentication.zip

#src
wget -c --progress=dot:giga https://github.com/cyberb/jellyfin-plugin-ldapauth/archive/refs/heads/master.zip -O ldap-src.zip
unzip ldap-src.zip

#custom binary    
#wget -c --progress=dot:giga https://github.com/cyberb/jellyfin-plugin-ldapauth/releases/download/v12-memberuid/LDAP-Auth.tar.gz
