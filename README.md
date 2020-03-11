# sck

A basic CLI tool to change individual settings in an SSH config file

```
$ sck --host github.com --parameter IdentityKey --value '~/.ssh/foo-bar-baz' --dry-run
New SSH Config:
# host-based configuration

Host github.com
  foo bar
  IdentityFile ~/.ssh/foo-bar-baz
```