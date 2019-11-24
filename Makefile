compile-abieos:
	cd abieos && git submodule update --init --recursive
	cd abieos && rm -rf build lib && mkdir build && cd build && cmake .. && make
	cd abieos && mkdir -p lib/darwin && cp build/libabieos.dylib lib/darwin

update-abieos:
	git submodule update --init --remote --rebase