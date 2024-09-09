# dockerbx

`dockerbx` is a Docker-based alternative to toolbx, designed for macOS users who want to create and manage isolated development environments easily. It provides a simple CLI interface to create, enter, and manage containers, allowing developers to maintain separate environments for different projects without cluttering their main system.

## Features

- Create isolated development environments using Docker containers
- Enter containers with a simple command, starting them if they're not running
- List all dockerbx containers
- Remove containers easily, with options for force removal and bulk operations
- Seamless integration with the host file system
- Based on Fedora images for a familiar Linux environment

## Installation

To install `dockerbx`, you can just download the binary from the release page, or build the latest development version by yourself (you need to have Go and Docker installed on your macOS system).

1. Clone the repository:
   ```
   git clone https://github.com/albertoperdomo2/dockerbx.git
   ```

2. Navigate to the project directory:
   ```
   cd dockerbx
   ```

3. Build the project:
   ```
   go build -o dockerbx cmd/dockerbx/main.go
   ```

4. Move the binary to a directory in your PATH:
   ```
   sudo mv dockerbx /usr/local/bin/
   ```

## Usage

### Initialize dockerbx

Before using dockerbx, run the init command:

```
dockerbx init
```

This will set up the necessary Docker images and configurations.

### Create a new container

```
dockerbx create [container_name]
```

If no name is provided, it will use the default name from the config file.

### Enter a container

```
dockerbx enter [container_name]
```

This will start the container if it's not running and give you an interactive shell.

### List containers

```
dockerbx list
```

Shows all dockerbx containers, their status, and creation date.

### Remove containers

```
dockerbx rm [container_name...]
```

Remove one or more containers. Use flags for additional options:
- `-f` or `--force`: Force removal of running containers
- `-a` or `--all`: Remove all dockerbx containers

### Run a command in a container

```
dockerbx run [container_name] [command]
```

Executes a command in the specified container without entering it.

### Update a container

```
dockerbx update [container_name]
```

Updates the container's base image and optionally updates packages within the container.
Use the `-p` or `--packages` flag to update packages.

## Configuration

The configuration file is located at `~/.config/dockerbx/dockerbx.yaml`. You can modify this file to change default settings.

## Development

To contribute to dockerbx:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgements

- Inspired by the [toolbx](https://github.com/containers/toolbox) project
- Built with [Docker](https://www.docker.com/) and [Go](https://golang.org/)
