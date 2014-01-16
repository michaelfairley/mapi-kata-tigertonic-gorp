set -x

cd mapi-kata
bash setup.sh
cd -

go build && ./mapi-kata-tigertonic-gorp config.json &
PID=$!

cd mapi-kata
bash run.sh
RETVAL=$?
cd -

kill $PID

exit $RETVAL
