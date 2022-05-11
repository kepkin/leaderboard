Test Data Generation
--------------------

The command bellow will generate log normal distributed data for structure initialisation

	``go run ./test/cmd/gen-test-data init --path ./test/data/insertion.csv``

To generate ammo for web API testing:

	``go run ./test/cmd/gen-test-data/ ammo --path ./test/data/ammo.txt``

Running yandex tank
-------------------

Start server

	``go build ./cmd/server
	./server``

Initialize with data

	``go run ./cmd/gen-test-data/ serverinit --path ./test/data/insertion.csv --endpoint http://localhost:8080/results``

Than run the yandex tank container

	``docker run -v $(pwd)/yandex-tank:/var/loadtest --net host -it direvius/yandex-tank``