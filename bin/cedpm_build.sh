#!/bin/bash

# The JSON file path
json_file='.cedpm-targets.json'

# Parse JSON and iterate over the instructions
jq -c '.targets[].instructions[]' "$json_file" | while read -r instruction; do
    # Extract command and input files from the instruction
    command=$(echo "$instruction" | jq -r '.command')
    inputs=$(echo "$instruction" | jq -r '.inputs[]')

    # Assuming the last file mentioned in the command is the output
    output=$(echo $command | awk '{print $NF}')

    # Check if output exists and compare timestamps
    should_execute=false
    if [ ! -f "$output" ]; then
        should_execute=true
    else
        for input in $inputs; do
            if [ "$input" -nt "$output" ]; then
                should_execute=true
                break
            fi
        done
    fi

    # Execute command if needed
    if $should_execute; then
        echo "Executing: $command"
        eval $command
    else
        echo "Skipping: $command (up-to-date)"
    fi
done
