mkdir -p /binaries/kotlin/1.6.10
cd /binaries/kotlin/1.6.10/
curl"https://github.com/JetBrains/kotlin/releases/download/v1.6.10/kotlin-native-linux-x86_64-1.6.10.tar.gz" \
  -o kotlin.tar.gz

tar xf kotlin.tar.gz --strip-components=1
rm kotlin.tar.gz