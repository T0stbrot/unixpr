# unixpr
Go Program that runs as a Probe Service for my Services, Linux Version

# Building it yourself
1. Install Go: https://go.dev/dl/
2. Download Source of latest release, unzip it, and navigate inside the foolder
4. Run `go build .`
5. Copy the `unixpr` file to `/usr/local/bin/unixpr`
7. Copy the `unixpr.service` file to `/etc/systemd/system/unixpr.service`
8. Run `chmod +x /usr/local/bin/unixpr`
9. Run `systemctl enable --now unixpr`
## Done, it is installed now
