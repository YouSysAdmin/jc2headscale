# jc2headscale

CLI tool for managing user groups in a Headscale policy file.

This tool generates a policy file based on a policy template,  
template is your policy file to which the specified groups with users from Jump Cloud will be added.

[![Stand with Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/banner2-direct.svg)](https://github.com/vshymanskyy/StandWithUkraine/blob/main/docs/README.md)

## Install

```shell
go install github.com/yousysadmin/jc2headscale/cmd/jc2headscale@latest
```

```shell
# By default install to $HOME/.bin dir
curl -L https://raw.githubusercontent.com/yousysadmin/jc2headscale/master/scripts/install.sh | bash
```

## Usage

```
Collects information about Jumpcloud groups, group members
and prepare a group list for Headscale policy.

Usage:
  jc2headscale [command]

Available Commands:
  prepare     Prepare policy

Flags:
      --input-policy string    Headscale/Tailscale policy file template (default "./policy.hjson")
      --jc-api-key string      The Jumpcloud API key (can use env var JC_API_KEY) (default "")
      --output-policy string   Headscale prepared policy file (default "./current.json")
      --strip-email-domain     Strip e-mail domain (default true)
```

The `--strip-email-domain` flag must be set eq to `oid.strip_email_domain` in your Headscale server config,
this flag determines whether it is necessary to trim the domain from the user's email or not, by default is `true`.

## Example

```sh
// Fill policy user groups from Jumpcloud
JC_API_KEY=0000000 jc2headscale prepare --input-policy=policy.hjson --output-policy=out.json

// Setup policy to Headscale
headscale policy set -f out.json
```

You can specify the Jumpcloud group in your Headscale policy file,
group name in the Jumcploud should be eq to a policy group name without the 'group:' prefix.

In the below example,  
group `not-jc` isn't present in the Jumpcloud, and as a result, this group stays as is,  
the `network-all` group present in the Jumpcloud and will be supplemented by users from Jumpcloud if users exist for
this
group and the `admin` user stays in the group,  
the `network-prod` group is present in the Jumpcloud and will be filled by users from Jumpcloud if users exist for this
group.

```jsonc
{
  "groups": {
    "group:not-jc": ["admin"],
    "group:network-all": ["admin"],
    "group:network-prod": []
  },
  ....
}
```
