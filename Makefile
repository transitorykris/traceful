all:
	docker build -t transitorykris/traceful .

push:
	docker push transitorykris/traceful

clean:
	docker rmi transitorykris/traceful
