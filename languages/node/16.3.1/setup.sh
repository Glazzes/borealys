curl "https://nodejs.org/dist/v16.13.1/node-v16.13.1-linux-x64.tar.xz" node.tar.gz
tar xf node.tar.gz --strip-components=1
rm node.tar.gz
export PATH=$PATH:$PWD/bin