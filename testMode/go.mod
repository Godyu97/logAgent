module logAgent/testMod

go 1.19

require logAgent/tailf v0.0.0

require (
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)

replace (
	logAgent/etcd => ../etcd
	logAgent/tailf => ../tailf
)
