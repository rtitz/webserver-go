# Webserver in Go

### How to install this in Linux (Systemd)
```bash
cp -v webserver-go.service /etc/systemd/system/
mkdir -p /opt/webserver-go/bin
cp -v bin/webserver-go_linux-arm64 /opt/webserver-go/bin/webserver-go
systemctl daemon-reload
systemctl --no-pager status webserver-go
systemctl enable --now webserver-go
systemctl --no-pager status webserver-go
```

### How to uninstall this in Linux (Systemd)
```bash
systemctl disable --now webserver-go
systemctl --no-pager status webserver-go
rm -f /etc/systemd/system/webserver-go.service
systemctl daemon-reload
rm -rf /opt/webserver-go/
```
