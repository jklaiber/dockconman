[![Container Build](https://github.com/jklaiber/dockconman/actions/workflows/container_build.yaml/badge.svg)](https://github.com/jklaiber/dockconman/actions/workflows/container_build.yaml)
# Docker Connection Manager
Simple connection manager written in Go an inspired from the awesome [ssh2docker](https://github.com/moul/ssh2docker/) project. 

## Usage
The basic idea of this connection manager was to build a tool which allows you to connect via SSH to a defined running container with a specific execution command. This allows you to let users access specific functions of a container without exposing more services which you don't want to expose.  
  
A simple example of a use case can be found below.
```
+-------------+     +----------+     +--------------------+     +-------------+
|user terminal| --> |dockconman| --> |container with virsh| --> |virsh console|
+-------------+     +----------+     +--------------------+     +-------------+
```

### Basic Usage 
Start the programm and define the container (on which the user should be redirected) and the command (which should be executed on the remote container).
```
$ dockconman -n <container_name> -c bash
```
Connect to the container via SSH and you will be redirected.
```
$ ssh localhost -p 2222
```
For additional settings you can use the `help` command.
```
$ dockconman help
NAME:
   dockconman - simple ssh portal to containers

USAGE:
   dockconman [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --container_name value, -n value   Target container [$DOCKCONMAN_CONTAINER]
   --command value, -c value          Execute command on target container (default: "bash") [$DOCKCONMAN_COMMAND]
   --user value, -u value             User for authenticating [$DOCKCONMAN_USER]
   --password value, -p value         Password for authenticating [$DOCKCONMAN_PASSWORD]
   --key-destination value, -k value  Host key destination which should be taken [$DOCKCONMAN_SSH_KEY_FILE]
   --default-key value, -d value      Disable automatic key generation (default: "false") [$DOCKCONMAN_DEFAULT_KEY]
   --banner value, -b value           Login banner (default: "#############################\n# Docker Connection Manager #\n#############################\n")
   --shell value, -s value            Default shell (default: "/bin/sh")
   --port value                       Binding port (default: "2222") [$DOCKCONMAN_PORT]
   --help, -h                         show help (default: false)
```

### Docker Usage
When you want to use the program within a Docker container, then you can use the following command below.
```
$ docker run \
-it \
-p 2222:2222 \
--env DOCKCONMAN_CONTAINER=2434307646eb \
--env DOCKCONMAN_COMMAND=bash \
--env DOCKCONMAN_SSH_KEY_FILE=/etc/rsa/id_rsa \
-v /var/run/docker.sock:/var/run/docker.sock \
-v /etc/rsa/id_rsa:/etc/rsa/id_rsa \
jklaiber/dockconman:latest
```

## Environment Variables
| Variable | Description |
|---|---|
|DOCKCONMAN_CONTAINER| Target container |
|DOCKCONMAN_COMMAND| Command which should be executed on target container |
|DOCKCONMAN_PORT|Port on which the ssh server should listen|
|DOCKCONMAN_SSH_KEY_FILE|Host key destination which is mounted from the host system (e.g `/etc/rsa/id_rsa`)|
|DOCKCONMAN_USER|User which should be used for authenticating. When empty no authentication is used.|
|DOCKCONMAN_PASSWORD|Password which is used besides the username for authenticating. When empty no authentication is used.|
|DOCKCONMAN_BANNER|Banner which should be displayed after a successful login|

## SSH Key Handling
The tool provides three different handling methods with the host (container) ssh key. 
1. It can use a default key (specified with: `-d true`).
2. It can generate an own one when starting the application (this is the default).
3. You can mount your own key in the container or give it as value (specified with: `-k /etc/id_rsa`).