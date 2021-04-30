#!/bin/bash

start_clients_and_wait_for_exit() {
	declare -a pids=()
	declare -a not_killed_pids=()
	echo Spawning
	for i in {1..30}
	do
		../../bin/alm-mqtt-module-example >> /tmp/log 2>&1 &
		p=$!
		pids+=(${p})
		echo ${p}
		sleep 0.1
	done

	sleep 3

	echo Killing
	i=0
	for pid in "${pids[@]}"; do
		i=$((i+1))
		if !((i % 2)); then
			echo "Killing client: ${pid}"
			kill ${pid}
		else
			echo "Preserving client: ${pid}"
			not_killed_pids+=(${pid})
		fi
	done

	echo "Waiting for alive examples to exit"
	wait
}

echo "Spawning clients and wait for their exit"
start_clients_and_wait_for_exit
