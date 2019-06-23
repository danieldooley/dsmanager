### Sync
```
scp -r ./* ddooley@dooley-server.local:~/go/src/github.com/danieldooley/dsmanager
```

### Run

Server must be run from `~/go/src/github.com/danieldooley/dsmanager/cmd/dsmanager`

It needs `sudo` access to be able run `rtcwake`

So can be run with:
```
sudo /usr/local/go/bin/go run *.go
```

### Service

The `dsmanager.service` file can be used to setup dsmanager as a service