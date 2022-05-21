## Test Data Generation

The command bellow will generate log normal distributed data for structure initialisation

	go run ./test/cmd/gen-test-data init --path ./test/data/insertion.csv

To generate ammo for web API testing:

	go run ./test/cmd/gen-test-data/ ammo --path ./test/data/ammo.txt


## Getting insertion rate

	go test ./test -bench='.' -v -timeout 300m

	cd test/
	python3 render-benchmark-graphs.py benchmark-report-2022-05-11T12\:38\:20+03\:00/

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