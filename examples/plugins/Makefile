PLUGINS ?=  $(patsubst %/,%,$(sort $(dir $(wildcard */))))

all: $(PLUGINS)

.PHONY: $(PLUGINS)
$(PLUGINS):
	cd $@ && \
		go generate && \
		go build -o ../libgstgo$@.so -buildmode c-shared .
	rm libgstgo$@.h

clean:
	rm -f *.so *.h