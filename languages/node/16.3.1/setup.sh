mkdir -p /binaries/node/16.3.1
cd /binaries/node/16.3.1/
curl "https://nodejs.org/dist/v16.13.1/node-v16.13.1-linux-x64.tar.xz" -o node.tar.xz
tar xf node.tar.xz --strip-components=1
rm node.tar.xz