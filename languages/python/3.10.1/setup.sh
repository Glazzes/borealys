mkdir "python"
cd "python"
curl "https://www.python.org/ftp/python/3.10.1/Python-3.10.1.tar.xz" -o python.tar.gz
tar xf python.tar.gz --strip-components=1
export PATH=$PATH:$PWD/python/bin