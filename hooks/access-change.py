import subprocess

print(subprocess.check_output('snap run jellyfin.access-change', shell=True))


