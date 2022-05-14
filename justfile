set windows-powershell := true
msg := "update"

push:
    go mod tidy
    git add .
    git commit -m '{{msg}}'
    git push

run:
    go build
    kill `cat /var/run/notification.pid` || true
    nohup ./notification >n.log 2>&1 &
