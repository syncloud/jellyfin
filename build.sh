#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}
apt update
apt install -y libltdl7 libnss3 wget unzip

VERSION=$1

BUILD_DIR=${DIR}/build/snap/app

docker ps -a -q --filter ancestor=app:syncloud --format="{{.ID}}" | xargs docker stop | xargs docker rm || true
docker rmi app:syncloud || true
docker build --build-arg VERSION=$VERSION -t app:syncloud .
docker run app:syncloud dotnet --help || true
docker create --name=app app:syncloud
mkdir -p ${BUILD_DIR}
echo $VERSION > $BUILD_DIR/app.version
cd ${BUILD_DIR}
docker export app -o app.tar
docker ps -a -q --filter ancestor=app:syncloud --format="{{.ID}}" | xargs docker stop | xargs docker rm || true
docker rmi app:syncloud || true
tar xf app.tar
rm -rf app.tar
#mkdir -p $BUILD_DIR/plugins/LDAP-Auth
#unzip ${DIR}/build/ldap-authentication.zip -d $BUILD_DIR/plugins/LDAP-Auth
#cp $DIR/build/jellyfin-plugin-ldapauth-memberuid/out/LDAP-Auth.dll $BUILD_DIR/plugins/LDAP-Auth
#cp $DIR/build/jellyfin-plugin-ldapauth-memberuid/out/Novell.Directory.Ldap.NETStandard.dll $BUILD_DIR/plugins/LDAP-Auth
#cp $DIR/config/jellyfin/meta.json $BUILD_DIR/plugins/LDAP-Auth
mkdir -p $BUILD_DIR/plugins
tar xf $DIR/build/LDAP-Auth.tar.gz -C $BUILD_DIR/plugins