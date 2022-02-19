#!/bin/bash

# chmod +x /root/noty/scripts/login-notify.sh
# add to /etc/pam.d/sshd
# session optional pam_exec.so seteuid /root/noty/scripts/login-notify.sh

addr=$1
agentid=$2
username=$3

if [ "${PAM_TYPE}" != "close_session" ]; then
    host="$(hostname)"
    subject="SSH Login: ${PAM_USER} from ${PAM_RHOST} on ${host} at $(date --rfc-3339=seconds)"
    curl -X POST -d '{"to_username":"'"${username}"'","content":"'"${subject}"'"}' http://${addr}/qiye-wechat/text-senders/${agentid}
fi
