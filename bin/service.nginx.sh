#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

timeout 100 /bin/bash -c 'until [ -S /var/snap/jellyfin/current/socket ]; do echo "waiting for socket"; sleep 1; done'
/bin/rm -f /var/snap/jellyfin/common/web.socket
exec ${DIR}/nginx/sbin/nginx -c /snap/jellyfin/current/config/nginx.conf -p ${DIR}/nginx -e stderr
