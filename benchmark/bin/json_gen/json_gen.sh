#!/bin/bash

##
# Bash script that generates a json payload for harbors.
#
# Author: Tiago Melo (tiagoharris@gmail.com)
##

# validate input parameters.
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <number_of_ports> <output_file>"
    exit 1
fi

NUM_PORTS=$1
OUTPUT_FILE=$2

# start message.
echo "Generating ${NUM_PORTS} ports in ${OUTPUT_FILE}..."

# write opening curly brace.
echo "{" > "$OUTPUT_FILE"

# generate ports.
for i in $(seq 1 "$NUM_PORTS"); do
    if [ "$i" -ne "$NUM_PORTS" ]; then
        echo "\"PORT${i}\": { \"name\": \"Port${i}\", \"city\": \"City${i}\", \"country\": \"Country${i}\", \"coordinates\": [${i}, ${i}], \"province\": \"Province${i}\", \"timezone\": \"UTC\", \"unlocs\": [\"PORT${i}\"], \"code\": \"12345\" }," >> "$OUTPUT_FILE"
    else
        echo "\"PORT${i}\": { \"name\": \"Port${i}\", \"city\": \"City${i}\", \"country\": \"Country${i}\", \"coordinates\": [${i}, ${i}], \"province\": \"Province${i}\", \"timezone\": \"UTC\", \"unlocs\": [\"PORT${i}\"], \"code\": \"12345\" }" >> "$OUTPUT_FILE"
    fi
done

# append closing curly brace.
echo "}" >> "$OUTPUT_FILE"

# end message.
echo "Generation completed. File saved as ${OUTPUT_FILE}"
