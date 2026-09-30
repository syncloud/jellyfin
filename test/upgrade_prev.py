import re

import pytest
import requests
from syncloudlib.integration.hosts import add_host_alias
from syncloudlib.http import wait_for_rest

from test import lib

TMP_DIR = '/tmp/syncloud'


@pytest.fixture(scope="session")
def module_setup(request, device, artifact_dir):
    def module_teardown():
        device.run_ssh('journalctl > {0}/prev.journalctl.log'.format(TMP_DIR), throw=False)
        device.scp_from_device('{0}/*.log'.format(TMP_DIR), artifact_dir)

    request.addfinalizer(module_teardown)


def test_start(module_setup, app, device_host, domain, device):
    add_host_alias(app, device_host, domain)
    device.activated()
    device.run_ssh('rm -rf {0}'.format(TMP_DIR), throw=False)
    device.run_ssh('mkdir {0}'.format(TMP_DIR), throw=False)


def test_install_store_version(device, app, app_domain):
    device.run_ssh('snap remove {0}'.format(app))
    device.run_ssh('snap install {0}'.format(app), retries=10)
    wait_for_rest(requests.session(), "https://{0}".format(app_domain), 200, 10)


def test_stamp_encoding_config(device, snap_data_dir):
    path = lib.encoding_path(snap_data_dir)
    before = lib.read_file(device, path)
    after = re.sub(r'<H264Crf>\d+</H264Crf>',
                   '<H264Crf>{0}</H264Crf>'.format(lib.CRF_MARKER), before)
    assert after != before, before

    lib.write_file(device, path, after)
    assert '<H264Crf>{0}</H264Crf>'.format(lib.CRF_MARKER) in lib.read_file(device, path)


def test_stamp_server_identity(device, app_domain):
    info = lib.server_info(app_domain)
    lib.write_file(device, lib.ID_FILE, info['Id'])
    lib.write_file(device, lib.VERSION_FILE, info['Version'])
