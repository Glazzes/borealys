mkdir go
cd go
curl "https://go.dev/dl/go1.17.5.linux-amd64.tar.gz" -o go.tar.gz
tar xf go.tar.gz --strip-components=1
export PATH=$PATH:$PWD/go/bin