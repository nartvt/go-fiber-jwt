

rm -rf ./bin/server

make build 

nohup ./bin/server > app.log&

tail -f app.log
