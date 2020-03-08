PWD := ${CURDIR}
DOCKER_IMAGE_NAME := "wwwlioss"
BASE_URL := "https://tamada.github.io/lioss"
THEME := dot

start:
	docker run --rm -d --name ${DOCKER_IMAGE_NAME} -p 1313:1313 -v "${PWD}":/src -v "$(PWD)"/../docs:/output -e HUGO_THEME=${THEME} -e HUGO_WATCH="true" -e HUGO_BASEURL=${BASE_URL} jojomi/hugo:latest

stop:
	docker stop ${DOCKER_IMAGE_NAME}

build:
	docker run --rm -d --name ${DOCKER_IMAGE_NAME} -v "${PWD}":/src -v ${PWD}/../docs:/output -e HUGO_THEME=${THEME} -e HUGO_BASEURL=${BASE_URL} jojomi/hugo:latest
