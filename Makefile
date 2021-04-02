GOBLOG_HOME = ~/goblog/src

.PHONY: mount
mount:
	sudo mount --bind $(shell pwd) $(GOBLOG_HOME)
