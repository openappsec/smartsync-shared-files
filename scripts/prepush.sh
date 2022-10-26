#!/bin/sh

version="0.1.0"               # Sets version variable

# Set Flags
# -----------------------------------
# Flags which can be overridden by user input.
# Default values are below
# -----------------------------------
gtest=0
gimports=0
glint=0
gvet=0
gmod=0
quiet=0
verbose=0
debug=0
args=()

GTEST_CMD="go test"
GIMPORTS_CMD="goimports -l -w"
GLINT_CMD="golint -set_exit_status"
GVET_CMD="go vet"
GMOD_CMD="go mod"


PASS=true
function mainScript() {

    if [[ $allFiles == 1 ]]; then
      GO_FILES=$(getAllFiles)
    else
    GO_FILES=$(stagedFiles)
    fi

    if ! [ -z "$GO_FILES" ]
    then
        if [[ $gimports == 1 ]]; then
            runCmdOver "$GIMPORTS_CMD" "$GO_FILES"
        fi

        if [[ $glint == 1 ]]; then
            runCmdOver "$GLINT_CMD" "$GO_FILES"
        fi

        if [[ $gvet == 1 ]]; then
            runCmdOver "$GVET_CMD" "$GO_FILES"
        fi
    fi

    if [[ $gtest == 1 ]]; then
        PKG_LIST=$(testPkgList)
        print "Running $GTEST_CMD..."
        printVerbose "Running on these packages: $PKG_LIST"
        runVerbose "$GTEST_CMD $PKG_LIST"
        if [[ $? != 0 ]]; then
            PASS=false
            print "$GTEST_CMD failed"
        fi
    fi

    if [[ $gmod == 1 ]]; then
        runGmod "$GMOD_CMD tidy -v"
        runGmod "$GMOD_CMD download"
        runGmod "$GMOD_CMD vendor"
    fi

    if ! $PASS; then
      print "Failure"
      exit 1
    else
      print "Success"
    fi

    exit 0
}

runCmdOver() {
    LOCAL_PASS=true
    print "Running $1..."
    printVerbose "Running over: $2"
    for ITEM in $2
    do
        runVerbose "$1 $ITEM"
        if [[ $? != 0 ]]; then
            LOCAL_PASS=false
        fi
    done
    if ! $LOCAL_PASS; then
        print "ERROR: $1 failed"
        PASS=false
    fi
}

runGmod() {
    print "Running $1..."
    runVerbose $1
    if [[ $? != 0 ]]; then
        PASS=false
        print "ERROR: $1 failed"
    fi
}

getAllFiles() {
    result=$(find 2> /dev/null | grep ".go$" | grep -v /mocks | grep -v /testdata )
    echo "$result"
}

stagedFiles() {
    result=$(git status -s 2> /dev/null | grep -v "D" | grep -v "R" | grep -v "C" | grep ".go$" | grep -v /mocks | grep -v /testdata | awk '{print $2}')
    echo "$result"
}

testPkgList() {
    result=$(go list ./... | grep -v /mocks | grep -v /test | grep -v /ingector | grep -v /cmd)
    echo "$result"
}

runVerbose() {
    if [[ $verbose == 1 ]] && [[ $quiet == 0 ]]; then
         $@ 2>&1 | sed 's/^/    /'
    fi
    $@ &> /dev/null
}

printVerbose() {
    if [[ $verbose == 1 ]] && [[ $quiet == 0 ]]; then
        echo $1
    fi
}

print() {
    if [[ $quiet == 0 ]]; then
        echo $1
    fi
}

# Print usage
usage() {
  echo "Usage: prepush.sh [OPTION]...

Golang prepush script

 Options:
  -c, --common      Run the script with 'common' values. Runs '$GLINT_CMD' and '$GIMPORTS_CMD' in verbose mode.
  -t, --gtest       Run '$GTEST_CMD'
  -i, --gimports    Run '$GIMPORTS_CMD'.
  -l, --glint       Run '$GLINT_CMD'.
  -e, --gvet        Run '$GVET_CMD'.
  -m, --gmod        Run '$GMOD_CMD tidy', '$GMOD_CMD download' and '$GMOD_CMD vendor' to update dependencies.
  -q, --quiet       Quiet (no output)
  -a  --all         Run on all go files, not hust staged files
  -v, --verbose     Output more information. (Items echoed to 'verbose')
  -d, --debug       Runs script in BASH debug mode (set -x)
  -h, --help        Display this help and exit
      --version     Output version information and exit
"
}

# Iterate over options breaking -ab into -a -b when needed and --foo=bar into
# --foo bar
optstring=h
unset options
while (($#)); do
  case $1 in
    # If option is of type -ab
    -[!-]?*)
      # Loop over each character starting with the second
      for ((i=1; i < ${#1}; i++)); do
        c=${1:i:1}

        # Add current char to options
        options+=("-$c")

        # If option takes a required argument, and it's not the last char make
        # the rest of the string its argument
        if [[ $optstring = *"$c:"* && ${1:i+1} ]]; then
          options+=("${1:i+1}")
          break
        fi
      done
      ;;

    # If option is of type --foo=bar
    --?*=*) options+=("${1%%=*}" "${1#*=}") ;;
    # add --endopts for --
    --) options+=(--endopts) ;;
    # Otherwise, nothing special
    *) options+=("$1") ;;
  esac
  shift
done
set -- "${options[@]}"
unset options

# Print help if no arguments were passed.
# Uncomment to force arguments when invoking the script
[[ $# -eq 0 ]] && set -- "--help"

# Read the options and set stuff
while [[ $1 = -?* ]]; do
  case $1 in
    -h|--help) usage >&2; exit 0 ;;
    --version) echo "$(basename $0) ${version}"; exit 0 ;;
    -c| --common) gimports=1; glint=1; gmod=0; gtest=0; gvet=0; verbose=1; allFiles=0 ;;
    -t|--gtest) gtest=1 ;;
    -i|--gimports) gimports=1 ;;
    -l|--glint) glint=1 ;;
    -e|--gvet) gvet=1 ;;
    -m|--gmod) gmod=1 ;;
    -a|--all) allFiles=1 ;;
    -v|--verbose) verbose=1 ;;
    -q|--quiet) quiet=1 ;;
    -d|--debug) debug=1;;
    --endopts) shift; break ;;
    *) echo "invalid option: '$1'."; exit 1 ;;
  esac
  shift
done

############## End Options and Usage ###################
# Run in debug mode, if set
if [ "${debug}" == "1" ]; then
  set -x
fi

# Bash will remember & return the highest exitcode in a chain of pipes.
set -o pipefail

# Run your script
mainScript

exit 0
