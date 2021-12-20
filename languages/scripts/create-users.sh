#!/bin/bash
for i in {1..100}
do
  currentUser="user$i"
  useradd -m $currentUser -d "/tmp/$currentUser"
  echo "$currentUser soft nproc 64" >> /etc/security/limits.conf
  echo "$currentUser hard nproc 64" >> /etc/security/limits.conf
  echo "$currentUser soft nofile 64" >> /etc/security/limits.conf
  echo "$currentUser hard nofile 64" >> /etc/security/limits.conf
  echo "$currentUser soft fsize 5120" >> /etc/security/limits.conf
  echo "$currentUser hard fsize 5120" >> /etc/security/limits.conf
  echo "$currentUser soft memlock 40960" >> /etc/security/limits.conf
  echo "$currentUser hard memlock 40960" >> /etc/security/limits.conf
done