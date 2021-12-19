mkdir node
cd node
curl "https://nodejs.org/dist/v16.13.1/node-v16.13.1-linux-x64.tar.xz" -o node.tar.gz
tar xf node.tar.gz --strip-components=1
export PATH=$PATH:$PWD/node/bin