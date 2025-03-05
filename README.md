# unixpr
Go Program that runs as a Probe Service for my Services, Linux Version

# Building it yourself
1. Install Go: https://go.dev/dl/
2. Download Source of latest release, unzip it, and navigate inside the foolder
4. Run `go build .`
5. Copy the `winpr` file to /usr/local/bin/winpr
7. Copy the `winpr.service` file to /etc/systemd/system/winpr.service
8. Run `systemctl enable --now winpr`
## Done, it is installed now
