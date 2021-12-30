mkdir -p /binaries/javascript/
cd /binaries/javascript/
curl "https://nodejs.org/dist/v16.13.1/node-v16.13.1-linux-x64.tar.xz" -o javascript.tar.xz
tar xf javascript.tar.xz --strip-components=1
rm javascript.tar.xz