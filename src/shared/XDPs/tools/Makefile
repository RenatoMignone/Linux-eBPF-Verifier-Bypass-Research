SUBDIRS := $(wildcard test_*/)

.PHONY: all clean $(SUBDIRS)

all:
	@for dir in $(SUBDIRS); do \
		echo ">>> Building in $$dir..."; \
		$(MAKE) -C $$dir; \
	done

clean:
	@for dir in $(SUBDIRS); do \
		echo ">>> Cleaning in $$dir..."; \
		$(MAKE) -C $$dir clean; \
	done

