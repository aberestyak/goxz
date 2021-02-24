# goxz

Tool for creating ssh reverse tunnel on remote host to local SOCKS server. Equivalent of:
`ssh -D <port> <user>@<address> “ssh -R <remote-address>:<remote-port>:<address>:<port> <remote-user>@<remote-host> -p <remote-port>”`

Specify following mandotory environment variables:

| ENV                   | Meaning                                                  |
| --------------------- | -------------------------------------------------------- |
| `SOCKS_HOST`          | SOCKS server address on local machine                    |
| `SOCKS_PORT`          | SOCKS server port on local machine                       |
| `REMOTE_SSH_HOST`     | Address of remote ssh server                             |
| `REMOTE_SSH_PORT`     | Port of remote ssh server                                |
| `REMOTE_SSH_USER`     | User to login on remote ssh server                       |
| `REMOTE_FORWARD_HOST` | Address on remote host to listen and forward connections |
| `REMOTE_FOWARD_PORT`  | Port on remote host to listen and forward connections    |

For enabling SOCKS server:

| ENV             | Meaning                               |
| --------------- | ------------------------------------- |
| `SOCKS_ENABLED` | Set `true` to enable SOCKS server     |
| `SOCKS_HOST`    | SOCKS server address on local machine |
| `SOCKS_PORT`    | SOCKS server port on local machine    |

Also you must specify at least one of environment variables to use for authentication on remote ssh server:

| ENV                 | Meaning                                                                                  |
| ------------------- | ---------------------------------------------------------------------------------------- |
| `SSH_KEY_PATH`      | Path to ssh private key. Use with `SSH_KEY_PASSPHRASE` env, if there is passhrase on key |
| `SSH_AUTH_SOCK`     | SSH-Agent socket address                                                                 |
| `SSH_USER_PASSWORD` | Password for user to login on remote ssh server                                          |

## TODO:

- [x] Dockerfile for build and runtime
- [ ] Implement tun2socks?
