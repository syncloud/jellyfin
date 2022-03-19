import subprocess

print(subprocess.check_output('snap run jellyfin.storage-change', shell=True))
