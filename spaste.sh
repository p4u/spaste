#!/bin/sh
SERVER="${SERVER:-paste.libremesh.org}"
PORT="${PORT:-80}"
HOST="$SERVER"
BODY="$(cat)"
BODY_LEN=$(echo -e "${BODY}\n" | wc -c)
URL="/add"
[ -n "$1" ] && URL="$URL/$1"
echo $SERVER | egrep -q "[A-z]" && {
  SERVER=$(ping -c1 -w1 -I lo $SERVER 2>/dev/null  | grep ^PING | cut -d"(" -f2 | cut -d")" -f1)
}
echo -ne "POST $URL HTTP/1.0
Host: $HOST
Content-Type:application/x-www-form-urlencoded
Content-Length: ${BODY_LEN}\r\n
${BODY}\r\n" | nc $SERVER $PORT | grep "^key=" | cut -d= -f2
