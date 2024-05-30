local name = "jellyfin";
local version = "10.9.3";
local browser = "firefox";
local platform = '22.02';
local selenium = '4.21.0-20240517';
local deployer = 'https://github.com/syncloud/store/releases/download/4/syncloud-release';


local build(arch, test_ui, dind) = [{
    kind: "pipeline",
    type: "docker",
    name: arch,
    platform: {
        os: "linux",
        arch: arch
    },
    steps: [
    {
        name: "version",
        image: "debian:buster-slim",
        commands: [
            "echo $DRONE_BUILD_NUMBER > version"
        ]
    },
    {
        name: "download",
        image: "debian:buster-slim",
        commands: [
            "./download.sh "
        ]
    },
    {
        name: "build",
        image: "docker:" + dind,
        commands: [
            "./build.sh " + version
        ],
        volumes: [
            {
                    name: "dockersock",
                    path: "/var/run"
                }
        ]
    },
    {
        name: "package python",
        image: "docker:" + dind,
        commands: [
            "./python/build.sh"
        ],
        volumes: [
            {
                    name: "dockersock",
                    path: "/var/run"
                }
        ]
    },
    {
        name: "package",
        image: "debian:buster-slim",
        commands: [
            "VERSION=$(cat version)",
            "./package.sh " + name + " $VERSION "
        ]
    },
    {
        name: "test",
        image: "python:3.8-slim-buster",
        commands: [
          "APP_ARCHIVE_PATH=$(realpath $(cat package.name))",
          "cd test",
          "./deps.sh",
          "py.test -x -s test.py --distro=buster --domain=buster.com --app-archive-path=$APP_ARCHIVE_PATH --device-host=" + name + ".buster.com --app=" + name + " --arch=" + arch
        ]
    }] + ( if test_ui then [
{
            name: "selenium",
            image: "selenium/standalone-" + browser + ":" + selenium,
            detach: true,
            environment: {
                SE_NODE_SESSION_TIMEOUT: "999999",
                START_XVFB: "true"
            },
               volumes: [{
                name: "shm",
                path: "/dev/shm"
            }],
            commands: [
                "cat /etc/hosts",
                "getent hosts " + name + ".buster.com | sed 's/" + name +".buster.com/auth.buster.com/g' | sudo tee -a /etc/hosts",
                "cat /etc/hosts",
                "/opt/bin/entry_point.sh"
            ]
         },

    {
        name: "selenium-video",
        image: "selenium/video:ffmpeg-4.3.1-20220208",
        detach: true,
        environment: {
                DISPLAY_CONTAINER_NAME: "selenium",
                FILE_NAME: "video.mkv"
            },
            volumes: [
            {
                name: "shm",
                path: "/dev/shm"
            },
            {
                name: "videos",
                path: "/videos"
            }
]
    },
    {
        name: "test-ui",
        image: "python:3.8-slim-buster",
        commands: [
          "apt-get update && apt-get install -y sshpass openssh-client libxml2-dev libxslt-dev build-essential libz-dev curl",
          "cd test",
          "pip install -r requirements.txt",
          "py.test -x -s ui.py --distro=buster --ui-mode=desktop --domain=buster.com --device-host=" + name + ".buster.com --app=" + name + " --browser=" + browser,
        ]
    }
     ] else [] ) +
   ( if arch == "amd64" then [
    {
        name: "test-upgrade",
        image: "python:3.8-slim-buster",
        commands: [
          "APP_ARCHIVE_PATH=$(realpath $(cat package.name))",
          "cd test",
          "./deps.sh",
          "py.test -x -s test-upgrade.py --distro=buster --ui-mode=desktop --domain=buster.com --app-archive-path=$APP_ARCHIVE_PATH --device-host=" + name + ".buster.com --app=" + name + " --browser=" + browser,
        ],
        privileged: true,
        volumes: [{
            name: "videos",
            path: "/videos"
        }]
    } ] else [] ) + [
    {
              name: 'upload',
              image: 'debian:buster-slim',
              environment: {
                AWS_ACCESS_KEY_ID: {
                  from_secret: 'AWS_ACCESS_KEY_ID',
                },
                AWS_SECRET_ACCESS_KEY: {
                  from_secret: 'AWS_SECRET_ACCESS_KEY',
                },
                SYNCLOUD_TOKEN: {
                  from_secret: 'SYNCLOUD_TOKEN',
                },
              },
              commands: [
                'PACKAGE=$(cat package.name)',
                'apt update && apt install -y wget',
                'wget ' + deployer + '-' + arch + ' -O release --progress=dot:giga',
                'chmod +x release',
                './release publish -f $PACKAGE -b $DRONE_BRANCH',
              ],
              when: {
                branch: ['stable', 'master'],
                event: ['push'],
              },
            },
            {
                  name: 'promote',
                  image: 'debian:buster-slim',
                  environment: {
                    AWS_ACCESS_KEY_ID: {
                      from_secret: 'AWS_ACCESS_KEY_ID',
                    },
                    AWS_SECRET_ACCESS_KEY: {
                      from_secret: 'AWS_SECRET_ACCESS_KEY',
                    },
                    SYNCLOUD_TOKEN: {
                      from_secret: 'SYNCLOUD_TOKEN',
                    },
                  },
                  commands: [
                    'apt update && apt install -y wget',
                    'wget ' + deployer + '-' + arch + ' -O release --progress=dot:giga',
                    'chmod +x release',
                    './release promote -n ' + name + ' -a $(dpkg --print-architecture)',
                  ],
                  when: {
                    branch: ['stable'],
                    event: ['push'],
                  },
                },
    {
        name: "artifact",
        image: "appleboy/drone-scp:1.6.4",
        settings: {
            host: {
                from_secret: "artifact_host"
            },
            username: "artifact",
            key: {
                from_secret: "artifact_key"
            },
            timeout: "2m",
            command_timeout: "2m",
            target: "/home/artifact/repo/" + name + "/${DRONE_BUILD_NUMBER}-" + arch,
            source: [
                "artifact/*"
            ],
            privileged: true,
            strip_components: 1,
            volumes: [
               {
                    name: "videos",
                    path: "/drone/src/artifact/videos"
                }
            ]
        },
        when: {
          status: [ "failure", "success" ]
        }
    }
    ],
    trigger: {
      event: [
        "push",
        "pull_request"
      ]
    },
    services:  [
{
                name: "docker",
                image: "docker:" + dind,
                privileged: true,
                volumes: [
                    {
                        name: "dockersock",
                        path: "/var/run"
                    }
                ]
            },
        {
            name: name + ".buster.com",
            image: "syncloud/platform-buster-" + arch + platform,
            privileged: true,
            volumes: [
                {
                    name: "dbus",
                    path: "/var/run/dbus"
                },
                {
                    name: "dev",
                    path: "/dev"
                }
            ]
        }
    ],
    volumes: [
        {
            name: "dbus",
            host: {
                path: "/var/run/dbus"
            }
        },
        {
            name: "dev",
            host: {
                path: "/dev"
            }
        },
        {
            name: "shm",
            temp: {}
        },
        {
            name: "videos",
            temp: {}
        },
        {
                name: "dockersock",
                temp: {}
            },
    ]
}];

build("amd64", true, "20.10.21-dind") +
build("arm64", false, "19.03.8-dind") +
build("arm", false, "19.03.8-dind")
