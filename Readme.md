Image Filter API:
----

Changes to be done:

1. Create .env file with name imagefilter.env and put variable names inside it

```

PROCESS_NAME = image_filter_server
LOG_LEVEL = info
ENVIRONMENT = dev
PORT = 8080
FIRESTORE_COLLECTION_NAME = <Collection Name>
FIRESTORE_PROJECT_ID = <Project Id>
GOOGLE_APPLICATION_CREDENTIALS= <Path to GCP Key Json File>

```

2. Download GCP Key Config Json from GCP console and rename it to `imagefilterapi-config.json`

3. In main.go , while initialising the module, change the input parameter to <Path to .env file>