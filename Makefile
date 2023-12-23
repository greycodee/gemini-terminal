BINARY_NAME=gemini-terminal
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME)

install: build
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)

clean:
	if [ -f $(BINARY_NAME) ] ; then rm $(BINARY_NAME) ; fi

.PHONY: build install clean