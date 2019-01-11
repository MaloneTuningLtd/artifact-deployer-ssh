# artifact-deployer-ssh

This is a pure Golang Drone >= 0.5 plugin to deploy binaries to a remote host.
For more information check back later.

## Docker
Build the Docker image by running:

```
docker build --rm=true -t maloneweb/artifact-deployer-ssh .
```

## Usage
Execute from the working directory (assuming you have an SSH server running on 127.0.0.1:22):

```
docker run --rm \
  -e PLUGIN_SSH_HOST="127.0.0.1" \
  -e PLUGIN_SSH_USER="web" \
  -e PLUGIN_SSH_PRIVATEKEY="$(cat some-private-key)" \
  -e PLUGIN_BINARY=$(pwd)/my-bin \
  -e PLUGIN_TAG=$DRONE_TAG \
  -e PLUGIN_PROJECT_PATH="/home/web/apps/bin-project" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  maloneweb/artifact-deployer-ssh
```
