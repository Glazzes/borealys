curl "https://go.dev/dl/go1.17.5.linux-amd64.tar.gz" -o go.tar.gz
tar xf go.tar.gz --strip-components=1
rm go.tar.gz
export PATH=$PATH:$PWD/bin