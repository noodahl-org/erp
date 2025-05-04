#!/bin/bash
KEY=$(cat ssh_key_base64.txt)
echo $KEY | base64 -d > temp_key
chmod 600 temp_key
ssh -i temp_key -o StrictHostKeyChecking=no caleb@192.168.1.146 "echo 'Connection successful'"
rm temp_key
