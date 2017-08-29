# Openstack Keystone plugin

This is a standalone backend plugin for use with [Hashicorp Vault](https://www.github.com/hashicorp/vault). This plugin provides the functionality to generate users in Openstack Keystone.



## Getting Started

This is a [Vault plugin](https://www.vaultproject.io/docs/internals/plugins.html)
and is meant to work with Vault. This guide assumes you have already installed Vault
and have a basic understanding of how Vault works.

Otherwise, first read this guide on how to [get started with Vault](https://www.vaultproject.io/intro/getting-started/install.html).

To learn specifically about how plugins work, see documentation on [Vault plugins](https://www.vaultproject.io/docs/internals/plugins.html).

### Build

- `go get github.com/parnurzeal/gorequest`
- `go get github.com/hashicorp/vault/plugins`
- `go get github.com/hashicorp/go-plugin`
- `go get github.com/fatih/structs`
- `go get github.com/google/gofuzz`
- `go build -o vault_keystone_plugin .``

### Installation

Build the plugin.

Put the plugin binary into a location of your choice. This directory
will be specified as the [`plugin_directory`](https://www.vaultproject.io/docs/configuration/index.html#plugin_directory)
in the Vault config used to start the server.

```json
...
plugin_directory = "path/to/plugin/directory"
...
```

Start a Vault server with this config file:
```sh
$ vault server -config=path/to/config.json ...
...
```

`sha256sum vault_keystone_plugin`

`vault write sys/plugins/catalog/vault_keystone_plugin sha_256="<SHA from the previous step>" command="keystone"`

`vault mount -path=keystone -plugin-name=vault_keystone_plugin plugin`

### Routes

## keystone/config/connection

CLI write / API POST - set connection configuration

Parameters:
-  `connection_url` : URL of your Keystone instance, formatted like `keystoneip:port/v3/`
-  `admin_auth_token` : admin user token

## keystone/users

CLI write / API POST
CLI read / API GET - generate new user

Parameters:
-  `name`
-  `default_project_id` (_optional_)
-  `domain_id` (_optional_)
-  `enabled` (_optional_)
-  `password` (_optional_)

## keystone/projects

CLI write / API POST
CLI read / API GET - generate new project

Parameters:
-  `name`
-  `is_domain` (_optional_)
-  `description` (_optional_)
-  `domain_id` (_optional_)
-  `enabled` (_optional_)
-  `parent_id` (_optional_)

## keystone/domains

CLI write / API POST
CLI read / API GET - generate new domain

Parameters:
-  `name`
-  `description` (_optional_)
-  `enabled` (_optional_)

## keystone/roles

CLI write / API POST
CLI read / API GET - generate new role

Parameters:
-  `name`
-  `domain_id` (_optional_)

## keystone/roles/*role*/groups/*group*/domains/*domain* action="grant"

CLI write / API POST - Assign role to group on domain

Parameters:
-  `domain_id`
-  `group_id`
-  `role_id`

## keystone/roles/*role*/users/*user*/domains/*domain* action="grant"

CLI write / API POST - Assign role to user on domain

Parameters:
-  `domain_id`
-  `user_id`
-  `role_id`

## keystone/roles/*role*/groups/*group*/projects/*project* action="grant"

CLI write / API POST - Assign role to group on project

Parameters:
-  `project_id`
-  `group_id`
-  `role_id`

## keystone/roles/*role*/users/*user*/projects/*project* action="grant"

CLI write / API POST - Assign role to user on project

Parameters:
-  `project_id`
-  `user_id`
-  `role_id`

### TODO:

- Keystone EC2 Extensions
- Credentials
- Groups
- Policies
- Regions
