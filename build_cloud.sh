cd cloudbuild

gcloud builds submit --config cloudbuild.car.yaml
gcloud builds submit --config cloudbuild.drivers.yaml
gcloud builds submit --config cloudbuild.events.yaml
