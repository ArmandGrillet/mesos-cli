# Mesos CLI for DC/OS

CLI for Mesos clusters running at the core of DC/OS.

Tested on DC/OS >= 1.10.

## Using the CLI as a DC/OS plugin

### On DC/OS >= 1.12

1. [Install the DC/OS CLI for your cluster and set it up](https://docs.mesosphere.com/latest/cli/install/).
2. Clone this repository and run `make zip` to generate the plugin.
3. Add the mesos-cli plugin to your CLI, e.g. using `dcos plugin add path/to/releases/mesos-cli.darwin.zip`.
4. Use the Mesos CLI by calling `dcos mesos <subcommand>`.

### On DC/OS >= 1.10

1. Install the cluster-agnostic DC/OS CLI by [downloading the latest release on GitHub](https://github.com/dcos/dcos-cli/releases).
2. Enable the plugin installation on cluster setup: `export DCOS_CLI_EXPERIMENTAL_AUTOINSTALL_PLUGINS=1`.
3. Clone this repository and run `make zip` to generate the plugin.
3. Set up the DC/OS CLI with your cluster using `dcos cluster setup`.
4. Add the mesos-cli plugin to your CLI, e.g. using `dcos plugin add path/to/releases/mesos-cli.darwin.zip`.
5. Use the Mesos CLI by calling `dcos mesos <subcommand>`.

## Using the CLI as a standalone

This requires having a DC/OS CLI working with your cluster and accessible using `dcos`.

1. Clone this repository and run `make`.
2. Set the correct environment variables by running `eval "$(./scripts/plugin_env.sh)"` (this uses `dcos`)
3. Use the CLI by using the binary, e.g. `build/darwin/bin/mesos-cli <subcommand>`.
