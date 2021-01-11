PLUGINS ?=  $(patsubst %/,%,$(sort $(dir $(wildcard */))))

all: $(PLUGINS)

.PHONY: $(PLUGINS)
$(PLUGINS):
	cd $@ && \
		go generate && \
		go build -o ../libgst$@.so -buildmode c-shared .
	rm libgst$@.h

clean:
	rm -f *.so *.h