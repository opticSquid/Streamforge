#!/bin/bash
read -p "Stream server address: " rtmp_server
read -p "Stream key: " stream_key
echo "Streaming to $rtmp_server/$stream_key"
ffmpeg -f v4l2 -framerate 30 -video_size 640x480 -i /dev/video0 -f pulse -i default -c:v libx264 -preset veryfast -tune zerolatency -b:v 2500k -g 60 -c:a aac -b:a 128k -f flv "$rtmp_server/$stream_key"
