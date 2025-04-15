.PHONEY: all docker-build docker-run
NAME := traffic_sniffer

all: docker-run

docker-build:
	docker build -t $(NAME) .

docker-run: docker-build
	docker run -d --privileged --net=host -v ./sniffing_logs:/logs $(NAME)
