#!/usr/bin/env bash
#
# Bash script to run the tests in the application and generate a coverage report
#
# This is done as a script because not all of the options are available in the Asure Pipeline Go tasks

# Define varibles
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
OUTPUT_DIR="${SCRIPT_DIR}/../build/tests"
REPORT_OUTPUT="report.xml"
COVERAGE_OUTPUT="coverage.xml"

# Analyse the arguments and configure the variables accordingly
while [[ $# -gt 0 ]]
do
    key ="$1"

    case $key in 
        -r|--report)
            REPORT_OUTPUT="$2"
        ;;

        -c|--coverage)
            COVERAGE_OUTPUT="$2"
        ;;
    esac
done

# Ensure that the output_dir exists
if [ ! -d ${OUTPUT_DIR} ]
then
    mkdir -p ${OUTPUT_DIR}
fi

# Install the necessary packages
go get github.com/jstemmer/go-junit-report
go get github.com/axw/gocov/gocov
go get github.com/AlekSi/gocov-xml

PATH+=":${HOME}/go/bin"

# Perform the tests int he correct directory
pushd internal

echo "Running tests: ${PWD}"

# - Unit tests
cmd="go test ./... -v | go-junit-report > ${OUTPUT_DIR}/${REPORT_OUTPUT}"
echo "Executing command: ${cmd}"
eval "${cmd}"

# check that the report file exists
if [ ! -f ${OUTPUT_DIR}/${REPORT_OUTPUT} ]
then
    echo "##vso[task.logissue type=error] Unable to find report file: ${OUTPUT_DIR}/${REPORT_OUTPUT}"
    exit 1
fi

# - Coverage
cmd="go test ./... -v -coverprofile=${OUTPUT_DIR}/cover.out"
echo "Executing command: ${cmd}"
eval "${cmd}"

ls -l ${OUTPUT_DIR}

if [ -f ${OUTPUT_DIR}/cover.out ]
then
    cmd="gocov convert ${OUTPUT_DIR}/cover.out | gocov-xml > ${OUTPUT_DIR}/${COVERAGE_OUTPUT}"
    echo "Executing command: ${cmd}"
    eval "${cmd}"
else
    echo "##vso[task.logissue type=error] Unable to find coverage file: ${OUTPUT_DIR}/cover.out"
    exit 1
fi

popd 