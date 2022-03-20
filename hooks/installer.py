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

APP_NAME = 'jellyfin'
USER_NAME = 'jellyfin'


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
        fs.makepath(join(self.data_dir, 'data', 'config'))
        fs.makepath(join(self.data_dir, 'data', 'plugins'))
        fs.makepath(join(self.data_dir, 'cache'))
        fs.makepath(join(self.data_dir, 'config'))
        
        fs.chownpath(self.data_dir, USER_NAME, recursive=True)
        fs.chownpath(self.common_dir, USER_NAME, recursive=True)

        self.prepare_storage()

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
        shutil.copy(join(self.snap_dir, 'config', 'system.xml'), join(self.data_dir, 'data', 'config', 'system.xml'))
        shutil.copytree(join(self.snap_dir, 'app', 'plugins', 'LDAP-Auth'), join(self.data_dir, 'data', 'plugins', 'LDAP-Auth'))

        app_storage_dir = storage.init_storage(APP_NAME, USER_NAME)
        with open(self.install_file, 'w') as f:
            f.write('installed\n')


    def prepare_storage(self):
        app_storage_dir = storage.init_storage(APP_NAME, USER_NAME)
        return app_storage_dir

