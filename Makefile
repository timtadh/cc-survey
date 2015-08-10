
assets.tar.gz: $(shell find ./assets/static ./assets/templates -type f)
	-rm -rf /tmp/assets
	mkdir -p /tmp/assets
	cp -r ./assets/static /tmp/assets/
	cp -r ./assets/templates /tmp/assets/
	tar czf assets.tar.gz -C /tmp assets/
	rm -rf /tmp/assets/

