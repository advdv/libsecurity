run:
	docker run -it -v /var/run/docker.sock:/var/run/docker.sock docksec

build:
	docker build -t docksec .