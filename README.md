# sck

A basic CLI tool to change individual settings in an SSH config file

```
$ sck host -h github.com --param IdentityKey --value '~/.ssh/foo-bar-baz' --dry-run
New SSH Config:

# global configuration

# host-based configuration

Host github.com
  IdentityKey ~/.ssh/foo-bar-baz
```