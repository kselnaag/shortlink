#!/usr/bin/env bash

echo -e "\n>>_ServerStart_<<"
./bin/shortlink 1>/dev/null &
SERVERPID=$!
if [[ $? -gt 0 ]]; then exit 1; fi

echo -e "\n>>_HealthCheck_<<"
sleep 2
RESPCODE=`curl -i http://localhost:8080/check/ping 2>/dev/null | head -n 1 | cut -d$' ' -f2`
if [ "$RESPCODE" != "200" ]; then 
    exit 1;
fi

echo -e "\n>>_CloseServer_<<\n"
kill $SERVERPID