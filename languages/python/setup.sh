mkdir -p /binaries/python/
cd /binaries/python/
curl "https://www.python.org/ftp/python/3.10.1/Python-3.10.1.tar.xz" -o python.tar.xz
tar xf python.tar.xz --strip-components=1
rm python.tar.xz

./configure
make install