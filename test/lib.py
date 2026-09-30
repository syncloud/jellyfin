import base64

import requests

CRF_MARKER = '19'
ID_FILE = '/tmp/syncloud/before-id'
VERSION_FILE = '/tmp/syncloud/before-version'


def server_info(app_domain):
    response = requests.get('https://{0}/System/Info/Public'.format(app_domain), verify=False)
    assert response.status_code == 200, response.text
    return response.json()


def read_file(device, path):
    output = device.run_ssh('base64 -w0 {0}'.format(path))
    lines = [line.strip() for line in output.splitlines() if line.strip()]
    assert lines, 'nothing read from {0}: {1}'.format(path, output)
    return base64.b64decode(lines[-1]).decode()


def write_file(device, path, content):
    blob = base64.b64encode(content.encode()).decode()
    device.run_ssh('echo {0} | base64 -d > {1}'.format(blob, path))


def encoding_xml(device, snap_data_dir):
    return read_file(device, encoding_path(snap_data_dir))


def encoding_path(snap_data_dir):
    return '{0}/config/encoding.xml'.format(snap_data_dir)
