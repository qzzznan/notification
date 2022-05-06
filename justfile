set windows-powershell := true
msg := "update"

push:
    git add .
    git commit -m '{{msg}}'
    git push

run:
    go build
    kill `cat /var/run/notification.pid` || true
    nohup ./notification --port 8006 >n.log 2>&1 &