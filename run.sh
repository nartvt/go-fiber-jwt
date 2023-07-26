

echo "REMOVE OLD BIN"
rm -rf ./bin/server
sleep 2

echo "BINARY BUILD ---------------->"
make build 
sleep 2
echo "DONE"

echo "KILL SERVICE------------->"
pkill -f "./bin/go_fiber_jwt_server"
sleep 2
echo "DONE"

echo "run SERVER"
nohup ./bin/go_fiber_jwt_server > app.log&
sleep 2

tail -f app.log
