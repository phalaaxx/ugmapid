ugmapid
===

The ugmapid is a tool for changing uid and gid of files and directories on a filesystem.
Its purpose is to make it easy to transfer unprivileged containers from one subuid and subgid range to another, like when transferring containers from one user account to another or between host servers.


installation
---

Compilation from sources requires golang 1.5+ compiler:

	go get https://github.com/phalaaxx/ugmapid


usage
---

ugmapid takes two arguments - directory and offset.
The directory argument specifies the root directory of the container and offset specifies the container's root uid/gid on the host server.
Make sure that the container uid and gid ranges are specified in /etc/subuid and /etc/subgid accordingly.

For example:

	ugmapid -directory /var/lib/lxc/container/rootfs -offset 100000

Note: uid and gid of directory before mapping are used to determine initial root uid and gid values.
For example, if /var/lib/lxc/container/rootfs is owned by root (uid=0 and gid=0) and offset is configured to 100000, after mapping root uid and gid will be set to 100000 and all other users will be mapped accordingly.

To convert unprivileged container to a privileged one it is possible to map it with offset = 0:

	ugmapid -directory /var/lib/lxc/containe/rootfs -offset 0
