PWD := ${CURDIR}
DOCKER_IMAGE_NAME := "wwwlioss"
BASE_URL := "https://tamada.github.io/lioss"
THEME := cayman-hugo-theme
JOJOMI_VERSION := 0.69.0

build:
	docker run --rm --name ${DOCKER_IMAGE_NAME} -v "${PWD}":/src -v ${PWD}/public:/output -e HUGO_THEME=${THEME} -e HUGO_BASEURL=${BASE_URL} jojomi/hugo:${JOJOMI_VERSION}

start:
	docker run --rm -it --name ${DOCKER_IMAGE_NAME} -p 1313:1313 -v "${PWD}":/src -v "$(PWD)"/public:/output -e HUGO_THEME=${THEME} -e HUGO_WATCH="true" -e HUGO_BASEURL=${BASE_URL} jojomi/hugo:${JOJOMI_VERSION}

stop:
	docker stop ${DOCKER_IMAGE_NAME}

