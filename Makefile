all:
	go build -o srcd cmd/runcore/*

clean:
	rm srcd
	rm -r ~/.srcd

run:
	./srcd account new
	./srcd --mine --miner.threads 4
