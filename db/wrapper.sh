#!/bin/sh

attempts=10
delay=5

for i in `seq $attempts`
do
	/go/bin/goose -env docker up && exit
	echo "Failed goose with exit code $?" 1>&2
	sleep $delay
done

echo "Failed to run goose migrations after $attempts" 1>&2
exit 1
