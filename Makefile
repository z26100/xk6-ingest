build:
        docker run --rm -it -e GOOS=windows -e GOARCH=amd64 -u "$(id -u):$(id -g)" -v "${PWD}/bin/windows:/xk6" grafana/xk6 build latest  --with github.com/z26100/xk6-ingest@latest
        docker run --rm -it -e GOOS=linux -e GOARCH=amd64 -u "$(id -u):$(id -g)" -v "${PWD}/bin/linux:/xk6" grafana/xk6 build latest  --with github.com/z26100/xk6-ingest@v0.3.0