#!/bin/bash

screen_updated="/data/screen/screen-updated"
if [ -f "$screen_updated" ]; then
    echo "$(date) screen-update: screen already up-to-date"
    exit 0
fi

tft_file=$(ls -t /tftpboot/*.tft | head -n1)

if [ -z "$tft_file" ]; then
    echo "$(date) screen-update: no file available"
    exit 0
fi

if [ ! -f "$tft_file" ]; then
    echo "$(date) screen-update: file not available"
    exit 0
fi

screen-updater -i $tft_file
res=$?
if [ "$res" != "0" ]; then
    echo "$(date) screen-update: exit with failure $res"
    exit 0
fi

$(touch "$screen_updated")
echo "$(date) screen-update: exit with success $res"
exit 0
