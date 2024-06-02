import json
import logging
import time
import uuid
from os import path
from os.path import join
from subprocess import check_output, CalledProcessError
import shutil
import requests
from syncloudlib import fs, linux, gen, logger
from syncloudlib.application import urls, storage
from syncloudlib.http import wait_for_rest
import re
import requests_unixsocket

APP_NAME = 'jellyfin'
USER_NAME = 'jellyfin'

SOCKET_FILE = '/var/snap/jellyfin/current/socket'
SOCKET = 'http+unix://{0}'.format(SOCKET_FILE.replace('/', '%2F'))


class Installer:
    def __init__(self):
        if not logger.factory_instance:
            logger.init(logging.DEBUG, True)

        self.log = logger.get_logger('jellyfin')
        self.snap_dir = '/snap/jellyfin/current'
        self.data_dir = '/var/snap/jellyfin/current'
        self.common_dir = '/var/snap/jellyfin/common'
        self.app_url = urls.get_app_url(APP_NAME)
        self.install_file = join(self.common_dir, 'installed')
         
    def pre_refresh(self):
        self.log.info('pre refresh')

    def post_refresh(self):
        self.log.info('post refresh')
        #self.init_config()

    def install(self):
        self.log.info('install')
        self.init_config()

    def init_config(self):
        linux.useradd(USER_NAME)

        log_dir = join(self.common_dir, 'log')
        fs.makepath(log_dir)
        fs.makepath(join(self.data_dir, 'nginx'))
        fs.makepath(join(self.data_dir, 'data'))
        
        fs.makepath(join(self.data_dir, 'data', 'plugins'))
        fs.makepath(join(self.data_dir, 'cache'))
        fs.makepath(join(self.data_dir, 'config'))
        shutil.copytree(join(self.snap_dir, 'app', 'plugins', 'LDAP-Auth'), join(self.data_dir, 'data', 'plugins', 'LDAP-Auth'))
        self.refresh_config()
        self.prepare_storage()

    def refresh_config(self):
        variables = {
            'domain': urls.get_app_domain_name(APP_NAME),
            'local_ipv4': self.local_ipv4(),
            'ipv6': self.ipv6()
        }
        gen.generate_files(join(self.snap_dir, 'config', 'jellyfin', 'config'), join(self.data_dir, 'config'), variables)
        shutil.copytree(join(self.snap_dir, 'config', 'jellyfin', 'plugins'), join(self.data_dir, 'data', 'plugins'), dirs_exist_ok=True)
        fs.chownpath(self.data_dir, USER_NAME, recursive=True)
        fs.chownpath(self.common_dir, USER_NAME, recursive=True)

    def local_ipv4(self):
        try:
            return check_output("/snap/platform/current/bin/cli ipv4", shell=True).decode().strip()
        except CalledProcessError as e:
            return 'localhost'

    def ipv6(self):
        try:
            return check_output("/snap/platform/current/bin/cli ipv6", shell=True).decode().strip()
        except CalledProcessError as e:
            return 'localhost'

    def configure(self):
        self.log.info('configure')
        if path.isfile(self.install_file):
            self._upgrade()
        else:
            self._install()
     
    def _upgrade(self):
        self.log.info('configure upgrade')

    def _install(self):
        self.log.info('configure install') 
        app_storage_dir = storage.init_storage(APP_NAME, USER_NAME)
        session = requests_unixsocket.Session()
        wait_for_rest(session, "{0}/web".format(SOCKET), 200, 100)
        session.post("{0}/Startup/Complete".format(SOCKET))
        with open(self.install_file, 'w') as f:
            f.write('installed\n')

    def on_domain_change(self):
        self.refresh_config()

    def prepare_storage(self):
        return storage.init_storage(APP_NAME, USER_NAME)
