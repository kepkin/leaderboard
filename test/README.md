## Test Data Generation

The command bellow will generate log normal distributed data for structure initialisation

	go run ./test/cmd/gen-test-data init --path ./test/data/insertion.csv

To generate ammo for web API testing:

	go run ./test/cmd/gen-test-data/ ammo --path ./test/data/ammo.txt


## Getting insertion rate

	go test ./test  -bench='InsertionRate' -benchtime=1x -v -timeout 300m

	cd test/
	./render-benchmark-graphs.py plotBenchmarkFolder benchmark-report-2022-05-11T12\:38\:20+03\:00/

### Comparing with previous version

To compare results with previous version you can specify two directories like that:

	./render-benchmark-graphs.py plotBenchmarkFolder \
	  benchmark-report-2022-05-11T12\:38\:20+03\:00/ \
	  benchmark-report-2022-05-20T12\:38\:20+03\:00/

To compare two specific insertion rates

	./render-benchmark-graphs.py cmpInsertionCSV \
	  benchmark-report-2022-05-22T16\:45\:56+03\:00/insertion-btree-m0s0.25.csv \
	  benchmark-report-2022-05-22T16\:45\:56+03\:00/insertion-btree-m0s1.csv \
	  benchmark-report-2022-05-22T16\:45\:56+03\:00/insertion-btree-m0s1-vs-m0s0.25



### Comparing with REDIS

Start Redis (with parameters not to write RDB during the test)
	
	docker run --rm -p 6379:6379 -it redis --save 3000 100000000



### Run short test while developing tests

	go test ./test -bench='.' -v -ldtestInitDuration 2s


## Running yandex tank

Start server

	go build ./cmd/server
	./server

Initialize with data

	go run ./cmd/gen-test-data/ serverinit --path ./test/data/insertion.csv --endpoint http://localhost:8080/results

Change ip address in 
Than run the yandex tank container

	docker run -v $(pwd)/yandex-tank:/var/loadtest --net host -it direvius/yandex-tank