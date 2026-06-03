#!/bin/bash
# complete -C 'mytask.sh -completion' mytask.sh
HERE=$(realpath .)

TOML="./mytask.toml"
while [ -z "${MYTASK_TOML}" ] && [ "$(pwd)" != "/" ]; do
        if [ -f $TOML ]; then
                MYTASK_TOML=$(realpath ${TOML})
        else
                cd $(dirname $(realpath .))
        fi
done


if [ -n "$MYTASK_TOML" ] ; then
	TASKDIR=$(toml-cli get ${MYTASK_TOML} mytask_dir 2> /dev/null)
	MYTASK_TASK_DIR=$(realpath ${TASKDIR:-.})
	COMPLETION=$(toml-cli get ${MYTASK_TOML} completion 2> /dev/null)

        cd $MYTASK_TASK_DIR
        if [ "$1" = "-completion" ]; then
                [ "${COMPLETION}" = "bash"] && go run . -toml $MYTASK_TOML -current -completion $HERE $@
        else
                go run . -toml $MYTASK_TOML -current $HERE $@
        fi
else
        echo mytask.toml file does not exists on any parent directory.
fi
cd $HERE

