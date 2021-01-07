PLUGINS ?=  $(patsubst %/,%,$(sort $(dir $(wildcard */))))

.PHONY: $(PLUGINS)
$(PLUGINS):
	cd $@ && \
		go build -o ../libgstgo$@.so -buildmode c-shared .
	rm libgstgo$@.h

clean:
	rm -f *.so *.h