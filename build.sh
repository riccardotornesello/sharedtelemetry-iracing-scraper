cd package
gcloud builds submit --config cloudbuild.yaml
cd ../qualify-results
gcloud builds submit --config cloudbuild.yaml
