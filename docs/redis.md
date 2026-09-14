# Redis

The Flynn Redis plugin provides Redis from the Ubuntu 24.04 package set in a
single process configuration. The data stored in this process is ephemeral and
is intended for caching and development use.

Install the plugin on the cluster with `flynn-host plugin install` (not the
user `flynn` CLI). After that, provision a database with `flynn resource add redis`.

## Usage

### Adding a server to an app

After the plugin is installed and you create an app, provision a server:

```text
flynn resource add redis
```

This provisions a Redis server as a Flynn app and configures your application
to connect to it.

### Connecting to the database

Provisioning the database will add a few environment variables to your app
release. `REDIS_HOST`, `REDIS_PORT`, and `REDIS_PASSWORD` provide connection
details for the database.

Flynn will also create the `REDIS_URL` environment variable which is utilized
by some libraries to configure connections.

### Connecting to a console

To connect to a console for the database, run `flynn redis redis-cli`. This does
not require the Redis client to be installed locally or firewall/security
changes, as it runs in a container on the Flynn cluster. The console uses
`REDIS_HOST` (`leader.<redis-app>.discoverd`), the same host as `REDIS_URL`.

### External access

An external route can be created that allows access to the database from
services that are not running on Flynn.

```text
flynn -a $(flynn env get FLYNN_REDIS) route add tcp --service $(flynn env get FLYNN_REDIS) --leader
```

This will provision a TCP port that always points at the primary instance.

For security reasons this port should be firewalled, and it should only be
accessed over the local network, VPN, or SSH tunnel.

## Safety

No safety or availability guarantees are currently provided for the Redis
appliance. Data loss and inconsistency is likely. Any data stored should be
treated as ephemeral and only used for caching, development, and testing.
