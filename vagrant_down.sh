#!/bin/bash
#Shutdown/destroy the dev env

function vagrant_down() {
    OUTPUT=`vagrant destroy -f 2>&1`
    #Check for: Your VM has become "inaccessible." Unfortunately, this is a critical error
    if [[ "$OUTPUT" == *"inaccessible"* ]]; then
        echo $OUTPUT
        echo "Apply fix https://stackoverflow.com/a/63087061/2032777"
        VBoxManage list vms |grep inaccessible |cut -d "{" -f2 |cut -d "}" -f1 |xargs -L1 VBoxManage unregistervm
    fi
    exit 0
}

vagrant_down
