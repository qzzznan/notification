#!/usr/bin/env bash

echo "
/var/log/notification/notification*.log {
        weekly
        missingok
        rotate 12
        compress
        notifempty
        postrotate
            kill -USR1 \`cat /var/run/notification.pid\`
        endscript
}
" >/etc/logrotate.d/notification
chmod 0644 /etc/logrotate.d/notification
#/usr/sbin/logrotate -d -f /etc/logrotate.d/notification
