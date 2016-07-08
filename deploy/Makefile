VERSION=v1.1.0

all:
		cd ..; make
		cp ../sigmactler .
		docker build -t sigmas/sigmactler:${VERSION} .
		echo ${VERSION}
		#docker push sigmas/sigmactler:${VERSION}

.PHONY: all
