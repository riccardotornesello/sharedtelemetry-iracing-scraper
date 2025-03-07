#!/bin/bash

set -e


docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/results-front:latest" --file docker/Dockerfile.results-front .
# docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/api:latest" --file docker/Dockerfile.api .
# docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/leagues-parser:latest" --file docker/Dockerfile.leagues-parser .
# docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/season-parser:latest" --file docker/Dockerfile.season-parser .
# docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/sessions-downloader:latest" --file docker/Dockerfile.sessions-downloader .
# docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/events-models:latest" --file docker/Dockerfile.events-models .

# docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/drivers-downloader:latest" --file docker/Dockerfile.drivers-downloader .
# docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/drivers-models:latest" --file docker/Dockerfile.drivers-models .

# docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/cars-downloader:latest" --file docker/Dockerfile.cars-downloader .
# docker build -t "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/cars-models:latest" --file docker/Dockerfile.cars-models .


docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/results-front:latest"
# docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/api:latest"
# docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/leagues-parser:latest"
# docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/season-parser:latest"
# docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/sessions-downloader:latest"
# docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/events-models:latest"

# docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/drivers-downloader:latest"
# docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/drivers-models:latest"

# docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/cars-downloader:latest"
# docker push "europe-west1-docker.pkg.dev/sharedtelemetryapp/sessions-downloader/cars-models:latest"