#!/bin/bash

TMP_DIR="/tmp/henge/cmd-tests"
TMP_STDOUT="${TMP_DIR}"/test-stdout
TMP_STDERR="${TMP_DIR}"/test-stderr

HENGE_ROOT=$(dirname "${BASH_SOURCE}")/../..

# in order to store stderr and stdout in different places, we need to store them into files
# this is why we need tmp directory
function initTmp(){
    mkdir -p ${TMP_DIR}
    rm -f ${TMP_STDOUT} ${TMP_STDERR}
}


# Utility function compare actual henge output with expected output
function compareOutput(){
    local target=${1}
    local dockerComposeFile=${2}
    local expectedOutput=${3}
    henge ${target} -f ${dockerComposeFile} > ${TMP_DIR}/compareOutput

    # skip ref,secret,uri field  - they are different every time
    diff --suppress-common-lines -y ${TMP_DIR}/compareOutput ${expectedOutput} \
        | grep -vE "ref\:|secret|uri\:" | tee ${TMP_STDOUT}
    if [[ "$(cat ${TMP_STDOUT} | wc -l)" -ne "0" ]]; then
        return 1
    else
        return 0
    fi
}

function checkDefaultFiles(){
    touch $2
    henge $1

    local exit_code=$?
    rm -f $2

    if [[ "${exit_code}" -eq "0" ]]; then
        return 0
    else
        return 1
    fi
}

# Run all "test_*" function from this file.
function runTests() {
    local failedTests=""
    #get all function names that are begginig with "test_"
    for testFce in $(declare -F | cut -f 3 -d ' ' | grep  -E "^test_");do
        echo "* Running ${testFce}"
        eval ${testFce} 2>$TMP_STDERR >$TMP_STDOUT
        local exit_code=$?
        if [[ $exit_code -ne 0 ]]; then
            failedTests="${failedTests} ${testFce}"
            echo "  TEST FAILED"
            echo "  stderr:"
            cat $TMP_STDERR
            echo ""
            echo "  stdout:"
            cat $TMP_STDOUT
        else
            echo "  test OK"
        fi
    done

    if [[ "$(echo $failedTests | wc -w)" -ne "0" ]]; then
        echo ""
        echo "FAILED TESTS: $failedTests"
        return 1
    else
        return 0
    fi
}


# regular henge run, verify right exit code
function test_exitCodeSuccess() {
   henge openshift -f ${HENGE_ROOT}/test/fixtures/complex/docker-compose.yml
    local exit_code=$?
    if [[ "${exit_code}" -eq "0" ]]; then
        return 0
    else
        return 1
    fi
}

# test right exit code when compose file doesn't exist
function test_fileNotExist(){
    henge openshift -f nonexiting_file
    local exit_code=$?
    if [[ "${exit_code}" -ne "0" ]]; then
        return 0
    else
        return 1
    fi
}

# test right exit for not supported target
function test_targetNotSupported(){
    henge nonexisting -f ${HENGE_ROOT}/test/fixtures/complex/docker-compose.yml
    local exit_code=$?
    if [[ "${exit_code}" -ne "0" ]]; then
        return 0
    else
        return 1
    fi
}

# test right exit code when non existing target
function test_nonExistigtarget(){
    henge nonexisting -f ${HENGE_ROOT}/test/fixtures/complex/docker-compose.yml
    local exit_code=$?
    if [[ "${exit_code}" -ne "0" ]]; then
        return 0
    else
        return 1
    fi
}


# check conversion for complex app to OpenShift
function test_complexOpenshift(){
    compareOutput "openshift" "${HENGE_ROOT}/test/fixtures/complex/docker-compose.yml" "${HENGE_ROOT}/test/fixtures/complex/docker-compose.converted.yaml"
    return $?
}


# check conversion for wordpress app to Kubernetes
function test_wordpressKubernetes(){
    compareOutput "kubernetes" "${HENGE_ROOT}/test/fixtures/wordpress/docker-compose.yml" "${HENGE_ROOT}/test/fixtures/wordpress/docker-compose.k8s.converted.yml"
    return $?
}


function test_DefaultDockerComposeYml(){
    checkDefaultFiles kubernetes docker-compose.yml
    checkDefaultFiles openshift docker-compose.yml
}


function test_DefaultDockerComposeYaml(){
    checkDefaultFiles kubernetes docker-compose.yaml
    checkDefaultFiles openshift docker-compose.yaml
}


initTmp
runTests
