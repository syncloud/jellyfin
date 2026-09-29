import re

import pytest
import requests
from subprocess import check_output
from syncloudlib.integration.hosts import add_host_alias
from syncloudlib.integration.installer import local_install
from syncloudlib.http import wait_for_rest

from test import lib

TMP_DIR = '/tmp/syncloud'


@pytest.fixture(scope="session")
def before(device):
    values = {
        'id': lib.read_file(device, lib.ID_FILE).strip(),
        'version': lib.read_file(device, lib.VERSION_FILE).strip(),
    }
    assert re.match(r'^[0-9a-f]{32}$', values['id']), values
    assert re.match(r'^[0-9][0-9.]*$', values['version']), values
    return values


@pytest.fixture(scope="session")
def module_setup(request, device, artifact_dir):
    def module_teardown():
        device.run_ssh('journalctl > {0}/refresh.journalctl.log'.format(TMP_DIR), throw=False)
        device.run_ssh('cp -r /var/snap/jellyfin/current/data/data/* {0}/'.format(TMP_DIR), throw=False)
        device.run_ssh('cp /var/snap/jellyfin/current/config/encoding.xml {0}/'.format(TMP_DIR), throw=False)
        device.scp_from_device('{0}/*'.format(TMP_DIR), artifact_dir)
        check_output('chmod -R a+r {0}'.format(artifact_dir), shell=True)

    request.addfinalizer(module_teardown)


def test_start(module_setup, app, device_host, domain, device):
    add_host_alias(app, device_host, domain)
    device.activated()


def test_upgrade(device_host, device_password, app_archive_path, app_domain):
    local_install(device_host, device_password, app_archive_path)
    wait_for_rest(requests.session(), "https://{0}".format(app_domain), 200, 20)


def test_server_reports_the_new_version(app_domain, before):
    version = lib.server_info(app_domain)['Version']
    assert version != before['version'], \
        'still on {0} after the refresh'.format(version)


def test_database_was_migrated_not_rebuilt(app_domain, before):
    info = lib.server_info(app_domain)
    assert info['Id'] == before['id'], \
        'server id changed from {0} to {1}, the database was rebuilt'.format(before['id'], info['Id'])
    assert info['StartupWizardCompleted'], info


def test_encoding_config_survives_the_upgrade(device, snap_data_dir):
    encoding = lib.encoding_xml(device, snap_data_dir)
    assert '<H264Crf>{0}</H264Crf>'.format(lib.CRF_MARKER) in encoding, encoding


def test_server_loads_every_config_file(device):
    journal = device.run_ssh('journalctl -u snap.jellyfin.server --no-pager')
    assert 'Error loading configuration' not in journal, journal


def test_ldap_plugin_is_loaded(device):
    journal = device.run_ssh('journalctl -u snap.jellyfin.server --no-pager')
    assert 'Loaded assembly LDAP-Auth' in journal, journal


def test_app_config_is_migrated(device, snap_data_dir, app_domain):
    network = device.run_ssh('cat {0}/config/network.xml'.format(snap_data_dir))
    assert app_domain in network, network

    logging = device.run_ssh('cat {0}/config/logging.default.json'.format(snap_data_dir))
    assert '"File"' not in logging, logging
