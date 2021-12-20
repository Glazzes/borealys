curl "https://www.python.org/ftp/python/3.10.1/Python-3.10.1.tar.xz" -o python.tar.gz
tar xf python.tar.gz --strip-components=1
rm python.tar.gz
export PATH=$PATH:$PWD/bin