#!/bin/bash

# creates runner users and limits available resources per user
for i in {1..3}
do
  currentUser="user$i"
  useradd -m $currentUser -G runners -d "/tmp/$currentUser"
  #echo "$currentUser soft nproc 128" >> /etc/security/limits.conf
  #echo "$currentUser hard nproc 128" >> /etc/security/limits.conf
  echo "$currentUser soft nofile 128" >> /etc/security/limits.conf
  echo "$currentUser hard nofile 128" >> /etc/security/limits.conf
  echo "$currentUser soft fsize 5120" >> /etc/security/limits.conf
  echo "$currentUser hard fsize 5120" >> /etc/security/limits.conf
  # echo "$currentUser soft memlock 40960" >> /etc/security/limits.conf
  # echo "$currentUser hard memlock 40960" >> /etc/security/limits.conf
done