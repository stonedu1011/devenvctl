#!/bin/bash -e

# shellcheck disable=SC1091
source "/usr/bin/versions.sh"

FILENAME="kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz"

url=$(curl --stderr /dev/null "https://www.apache.org/dyn/closer.cgi?path=/kafka/${KAFKA_VERSION}/${FILENAME}&as_json=1" | jq -r 'if .in_dist then .preferred+.path_info else "" end')

# Test to see if the suggested mirror has this version, currently pre 2.6.2 versions
if [[ ! $(curl -f -s -r 0-1 "${url}") ]]; then
    # do not appear to be actively mirrored. This may also be useful if closer.cgi is down.
    echo "Mirror does not have desired version, downloading direct from Apache"
    url="https://archive.apache.org/dist/kafka/${KAFKA_VERSION}/${FILENAME}"
fi

echo "Downloading Kafka from $url"
wget -q "${url}" -O "/tmp/${FILENAME}"
