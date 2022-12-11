#!/bin/bash
export day=$1
export dir="day${day}"

cp -R day_template "${dir}"
envsubst < day_template/go.mod > "${dir}/go.mod"

