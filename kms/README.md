Go services will be based on Microservices architecture. In essence, this is a wrapper service that will leverage Google Cloud Key Management Services to perform crypto services and deployed using Cloud Run in GCP

To Run and Test on local
========================
> Set env variables (Sample file found under kms/misc/envfile)
> Copy project Crdentialjson file as said in the env variable

Open a local terminal, go to kms folder and run service
>go run server.go

In another terminal try test
> cd kms/handlers
> go test .


To Run and test on local docker
===============================
Docker file got all the setting up for env variable and credential json file.
docker compose
> docker-compose up --build

In another terminal try test
> cd kms/handlers
> go test .


To Run and test on Cloud Run
============================
Cloudbuild.yaml file is used to build container image with three steps.
Build the container image, push image to container registry and finally deploy to cloud run.

