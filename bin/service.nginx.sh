#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [[ -z "$1" ]]; then
    echo "usage $0 [start|stop|reload]"
    exit 1
fi

case $1 in
start)
    timeout 100 /bin/bash -c 'until [ -S /var/snap/jellyfin/current/socket ]; do echo "waiting for socket"; sleep 1; done'
	  /bin/rm -f /var/snap/jellyfin/common/web.socket
    exec ${DIR}/nginx/sbin/nginx -c /snap/jellyfin/current/config/nginx.conf -p ${DIR}/nginx -g 'error_log stderr warn;'
    ;;
reload)
    ${DIR}/nginx/sbin/nginx -c /snap/jellyfin/current/config/nginx.conf -s reload -p ${DIR}/nginx
    ;;
stop)
    ${DIR}/nginx/sbin/nginx -c /snap/jellyfin/current/config/nginx.conf -s stop -p ${DIR}/nginx
    ;;
*)
    echo "not valid command"
    exit 1
    ;;
esac
