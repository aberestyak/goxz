# goxz

Tool for creating ssh reverse tunnel on remote host to local SOCKS server.
`server` mode is equivalent of:
`ssh -D <port> <user>@<address> “ssh -R <remote-address>:<remote-port>:<address>:<port> <remote-user>@<remote-host> -p <remote-port>”`

`client` mode is equivalent of:
`ssh -L <remote-address>:<remote-port>:<address>:<port> <remote-user>@<remote-host> -p <remote-port>”`

## Usage

```bash
goxz <server/client>
```

Can use

Specify following mandotory environment variables:

| ENV                   | Meaning                                                  |
| --------------------- | -------------------------------------------------------- |
| `LOCAL_HOST`          | address on local machine                                 |
| `LOCAL_PORT`          | port on local machine                                    |
| `REMOTE_SSH_HOST`     | Address of remote ssh server                             |
| `REMOTE_SSH_PORT`     | Port of remote ssh server                                |
| `REMOTE_SSH_USER`     | User to login on remote ssh server                       |
| `REMOTE_FORWARD_HOST` | Address on remote host to listen and forward connections |
| `REMOTE_FORWARD_PORT` | Port on remote host to listen and forward connections    |

> If remote host interface not configured ipv6, always listen on 0.0.0.0, no matter of REMOTE_FORWARD_HOST. [Issue](https://github.com/golang/go/issues/18806)

For enabling SOCKS server:

| ENV             | Meaning                           |
| --------------- | --------------------------------- |
| `SOCKS_ENABLED` | Set `true` to enable SOCKS server |


Also you must specify at least one of environment variables to use for authentication on remote ssh server:

| ENV                 | Meaning                                                                                  |
| ------------------- | ---------------------------------------------------------------------------------------- |
| `SSH_KEY_PATH`      | Path to ssh private key. Use with `SSH_KEY_PASSPHRASE` env, if there is passhrase on key |
| `SSH_AUTH_SOCK`     | SSH-Agent socket address                                                                 |
| `SSH_USER_PASSWORD` | Password for user to login on remote ssh server                                          |

## TODO:

- [x] Dockerfile for build and runtime
- [ ] Monitor goroutines (possible memory leaks)
- [ ] Implement tun2socks?
- [ ] Workaround for [Issue](https://github.com/golang/go/issues/18806)
