cd cloudbuild

gcloud builds submit --config ./cloudbuild/cloudbuild.cars.yaml
gcloud builds submit --config ./cloudbuild/cloudbuild.drivers.yaml
gcloud builds submit --config ./cloudbuild/cloudbuild.events.yaml
