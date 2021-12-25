mkdir -p /binaries/java/17
cd /binaries/java/17/
curl "https://cdn.azul.com/zulu/bin/zulu17.30.15-ca-jdk17.0.1-linux_x64.tar.gz" -o jdk-17.tar.gz
tar xf jdk-17.tar.gz --strip-components=1
rm jdk-17.tar.gz