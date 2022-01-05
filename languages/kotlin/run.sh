filename="$2"
compiledFilename=""

if [[ $filename =~ ^(/[a-z]+/[a-z]+/[a-z]+/).*$ ]]; then
    compiledFilename="${BASH_REMATCH[1]}code.jar"
else
  exit
fi

runuser -u "$1" -- /binaries/kotlin/kotlinc/bin/kotlinc "$2" -d "$compiledFilename"
runuser -u "$1" -- timeout -s KILL 4 /binaries/java/bin/java -jar "$compiledFilename"