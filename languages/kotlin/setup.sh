mkdir -p /binaries/kotlin/
cd /binaries/kotlin/
curl -L "https://github.com/JetBrains/kotlin/releases/download/v1.6.10/kotlin-compiler-1.6.10.zip" -o kotlin.zip
unzip kotlin.zip
rm kotlin.zip