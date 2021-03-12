# goxz

Tool for creating ssh reverse tunnel on remote host to local SOCKS server.
`server` mode is equivalent of:
`ssh -D <port> <user>@<address> “ssh -R <remote-address>:<remote-port>:<address>:<port> <remote-user>@<remote-host> -p <remote-port>”`

`client` mode is equivalent of:
`ssh -L <remote-address>:<remote-port>:<address>:<port> <remote-user>@<remote-host> -p <remote-port>”`

Forward UDP traffic over TCP:
- Open additional UDP endpoint (`CLIENT_HOST_UDP:CLIENT_PORT_UDP`) on client
- Open additional TCP-UDP endpoint (`SERVER_HOST_TCP-UDP:SERVER_PORT_TCP-UDP` **should be availiable from client!**) on server
- Client dials `SERVER_HOST_TCP-UDP:SERVER_PORT_TCP-UDP` and forward traffic from local UDP endpoint to this connection
- Server forwards traffic from TCP-UDP enpoint to specified UDP remote host (`REMOTE_HOST_UDP:REMOTE_PORT_UDP`) and send response over TCP back to client

## Usage

```bash
goxz <server/client>
```

### Specify following mandatory environment variables:

#### Server side:

| ENV                   | Meaning                                                                   |
| --------------------- | ------------------------------------------------------------------------- |
| `BASTION_HOST_SSH`    | Address of remote ssh server                                              |
| `BASTION_PORT_SSH`    | Port of remote ssh server                                                 |
| `BASTION_SSH_USER`    | User to login on remote ssh server                                        |
| `LOCAL_HOST_TCP`      | TCP address on local machine to forward TCP traffic                       |
| `LOCAL_PORT_TCP`      | TCP port on local machine to forward TCP traffic                          |
| `SERVER_HOST_TCP-UDP` | TCP address on local machine to forward TCP->UDP traffic (0.0.0.0 mainly) |
| `SERVER_PORT_TCP-UDP` | TCP port on local machine to forward TCP->UDP traffic                     |
| `REMOTE_HOST_UDP`     | Remote host UDP address (DNS server, for example)                         |
| `REMOTE_PORT_UDP`     | Remote host UDP address (DNS server, for example)                         |

> If remote host interface not configured ipv6, always listen on 0.0.0.0, no matter of REMOTE_FORWARD_HOST. [Issue](https://github.com/golang/go/issues/18806)


#### Client side

| ENV                   | Meaning                                                  |
| --------------------- | -------------------------------------------------------- |
| `BASTION_HOST_SSH`    | Address of remote ssh server                             |
| `BASTION_PORT_SSH`    | Port of remote ssh server                                |
| `BASTION_SSH_USER`    | User to login on remote ssh server                       |
| `LOCAL_HOST_TCP`      | TCP address on local machine to forward TCP traffic      |
| `LOCAL_PORT_TCP`      | TCP port on local machine to forward TCP traffic         |
| `SERVER_HOST_TCP-UDP` | Server TCP address to forward TCP->UDP traffic           |
| `SERVER_PORT_TCP-UDP` | Server TCP port to forward TCP->UDP traffic              |
| `CLIENT_HOST_UDP`     | UDP address on local machine to forward TCP->UDP traffic |
| `CLIENT_PORT_UDP`     | UDP port on local machine to forward TCP->UDP traffic    |


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
- [x] Implement UDP over TCP
- [ ] Workaround for [Issue](https://github.com/golang/go/issues/18806)
