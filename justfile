set windows-powershell := true
msg := "update"

build:
    go version

push:
    git add .
    git commit -m '{{msg}}'
    git push
