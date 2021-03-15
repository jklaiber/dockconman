# Docker Connection Manager
Simple connection manager written in Go an inspired from the awesome [ssh2docker](#https://github.com/moul/ssh2docker/) project. 

## Usage
The basic idea of this connection manager was to build a tool which allows you to connect via SSH to a defined running container with a specific execution command. This allows you to let users access specific functions of a container without exposing more services which you don't want to expose.  
  
A simple example of a use case can be found below.
```
ToDo Use Case
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
   --container_name value, -n value  Target container
   --command value, -c value         Execute command on target container (default: "bash")
   --banner value, -b value          Login banner (default: Docker Connection Manager)
   --shell value, -s value           Default shell (default: "/bin/sh")
   --port value, -p value            Binding port (default: ":2222")
   --help, -h                        show help (default: false)
```

### Docker Usage
When you want to use the program within a Docker container, then you can use the following command below.
```
$ 
```