import os
from os.path import dirname, join
from subprocess import check_output

import pytest
import requests
import time
from syncloudlib.integration.hosts import add_host_alias
from syncloudlib.integration.installer import local_install, wait_for_installer
from syncloudlib.http import wait_for_rest
DIR = dirname(__file__)
TMP_DIR = '/tmp/syncloud'


@pytest.fixture(scope="session")
def module_setup(device, request, data_dir, platform_data_dir, artifact_dir):
    def module_teardown():
        device.run_ssh('top -bn 1 -w 500 -c > {0}/top.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ps auxfw > {0}/ps.log'.format(TMP_DIR), throw=False)
        device.run_ssh('systemctl status snap.jellyfin.server > {0}/server.status.log'.format(TMP_DIR), throw=False)
        device.run_ssh('netstat -nlp > {0}/netstat.log'.format(TMP_DIR), throw=False)
        device.run_ssh('journalctl > {0}/journalctl.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/log/syslog {0}/syslog.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/log/messages {0}/messages.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /snap > {0}/snap.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /snap/jellyfin > {0}/snap.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /var/snap/jellyfin/current/> {0}/var.current.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la /var/snap/jellyfin/common/> {0}/var.common.ls.log'.format(TMP_DIR), throw=False)
        device.run_ssh('ls -la {0}/log > {1}/data.log.dir.ls.log'.format(data_dir, TMP_DIR), throw=False)
        device.run_ssh('df -h > {0}/df.log'.format(TMP_DIR), throw=False)
        device.run_ssh('df -h > {0}/df.log'.format(TMP_DIR), throw=False)

        app_log_dir = join(artifact_dir, 'log')
        os.mkdir(app_log_dir)
        device.scp_from_device('/var/snap/transmission/common/log/*.log', app_log_dir)
        device.scp_from_device('{0}/*'.format(TMP_DIR), app_log_dir)
        check_output('chmod -R a+r {0}'.format(artifact_dir), shell=True)

    request.addfinalizer(module_teardown)


def test_start(module_setup, device, app, domain, device_host):
    add_host_alias(app, device_host, domain)
    device.run_ssh('date', retries=10, throw=True)
    device.run_ssh('mkdir {0}'.format(TMP_DIR))
 

def test_activate_device(device):
    response = retry(device.activate_custom, 1)
    assert response.status_code == 200, response.text


def test_install(app_archive_path, device_host, device_password, app_domain):
    local_install(device_host, device_password, app_archive_path)
    wait_for_rest(requests.session(), "https://{0}".format(app_domain), 200, 10)
    

def test_reinstall(app_archive_path, device_host, device_password, app_domain):
    local_install(device_host, device_password, app_archive_path)
    wait_for_rest(requests.session(), "https://{0}".format(app_domain), 200, 10)


def test_ffmpeg(device, app_dir, data_dir):
    assert 'Universal media converter' in device.run_ssh('{0}/bin/ffmpeg.sh --help'.format(app_dir))


def test_ffprobe(device, app_dir, data_dir):
    assert 'streams analyzer' in device.run_ssh('{0}/bin/ffprobe.sh --help'.format(app_dir))


def test_storage_change(device, app_dir, data_dir):
    device.run_ssh('snap run jellyfin.storage-change > {1}/log/storage-change.log'.format(app_dir, data_dir),
                   throw=False)


def retry(method, retries=10):
    attempt = 0
    exception = None
    while attempt < retries:
        try:
            return method()
        except Exception as e:
            exception = e
            print('error (attempt {0}/{1}): {2}'.format(attempt + 1, retries, str(e)))
            time.sleep(5)
        attempt += 1
    raise exception


