# dockerbx Tutorial

This tutorial will guide you through using dockerbx to create and manage isolated development environments.

## Getting Started

1. Install dockerbx by downloading the latest release and adding it to your PATH.

2. Initialize dockerbx:
   ```
   dockerbx init
   ```
   This will set up the necessary Docker images and configurations.

## Creating Your First Environment

1. Create a new container:
   ```
   dockerbx create my-dev-env
   ```

2. Enter the container:
   ```
   dockerbx enter my-dev-env
   ```
   You're now in an isolated Fedora environment!

3. Install some tools:
   ```
   sudo dnf install -y python3 gcc
   ```

4. Exit the container:
   ```
   exit
   ```

## Running Commands Without Entering

You can run commands in your container without entering it:

```
dockerbx run my-dev-env python3 --version
```

## Updating Your Environment

To update the base image and packages:

```
dockerbx update my-dev-env --packages
```

## Cleaning Up

When you're done, you can remove the container:

```
dockerbx rm my-dev-env
```

## Advanced Usage

- Use the configuration file at `~/.config/dockerbx/dockerbx.yaml` to set default options.
- Create multiple environments for different projects.
- Use dockerbx in your development workflow to keep your host system clean.

Enjoy using dockerbx for your development needs!
