# DataStream to SQL Dataflow Template

The [DataStreamToSQL](src/main/java/com/google/cloud/teleport/v2/templates/DataStreamToSQL.java)
pipeline ingests data supplied by DataStream and writes the
data to a SQL Database via JDBC.
The supported databases include PostgreSQL and MySQL.

## Getting Started


### Requirements
* Java 11
* Maven
* DataStream stream is created and sending data to storage
* The SQL DB is accessible via Dataflow workers
* The SQL DB has schema matched to the Datastream DB.

### Building Template
This is a Flex Template meaning that the pipeline code will be containerized and the container will be
used to launch the Dataflow pipeline.

<!-- TODO(dhercher): Align everything below this line -->

#### Building Container Image
* Set environment variables.

```sh
export PROJECT=<my-project>
export IMAGE_NAME=datastream-to-postgres
export BUCKET_NAME=gs://<bucket-name>
export TARGET_GCR_IMAGE=gcr.io/${PROJECT}/${IMAGE_NAME}
export BASE_CONTAINER_IMAGE=gcr.io/dataflow-templates-base/java17-template-launcher-base
export BASE_CONTAINER_IMAGE_VERSION=latest
export APP_ROOT=/template/${IMAGE_NAME}
export DATAFLOW_JAVA_COMMAND_SPEC=${APP_ROOT}/resources/${IMAGE_NAME}-command-spec.json
export TEMPLATE_IMAGE_SPEC=${BUCKET_NAME}/images/${IMAGE_NAME}-image-spec.json

gcloud config set project ${PROJECT}
```

* Build and push image to Google Container Repository

```sh
mvn clean package \
-Dimage=${TARGET_GCR_IMAGE} \
-Dbase-container-image=${BASE_CONTAINER_IMAGE} \
-Dbase-container-image.version=${BASE_CONTAINER_IMAGE_VERSION} \
-Dapp-root=${APP_ROOT} \
-Dcommand-spec=${DATAFLOW_JAVA_COMMAND_SPEC} \
-am -pl ${IMAGE_NAME}
```

#### Creating Image Spec

* Create file in Cloud Storage with path to container image in Google Container Repository.
```sh
echo '{
    "image":"'${TARGET_GCR_IMAGE}'",
    "metadata":{"name":"PubSub CDC to BigQuery",
    "description":"Replicate Pub/Sub Data into BigQuery Tables",
    "parameters":[
        {
            "name":"inputSubscription",
            "label":"PubSub Subscription Name",
            "helpText":"Full subscription reference",
            "paramType":"TEXT"
        },
        {
            "name":"autoMapTables",
            "label":"Automatically add new BigQuery tables and columns as they appear",
            "helpText":"Automatically add new BigQuery tables and columns as they appear",
            "paramType":"TEXT",
            "isOptional":true
        },
        {
            "name":"outputDatasetTemplate",
            "label":"The BigQuery Dataset Name or column template",
            "helpText":"The BigQuery Dataset Name or column template",
            "paramType":"TEXT",
            "isOptional":true
        },
        {
            "name":"outputTableNameTemplate",
            "label":"The BigQuery Table Name or column template",
            "helpText":"The BigQuery Table Name or column template",
            "paramType":"TEXT",
            "isOptional":true
        },
        {
            "name":"outputTableSpec",
            "label":"DEPRECATED: Use outputDatasetTemplate AND outputTableNameTemplate",
            "helpText":"DEPRECATED: Use outputDatasetTemplate AND outputTableNameTemplate",
            "paramType":"TEXT",
            "isOptional":true
        },
        {"name":"outputDeadletterTable",
            "label":"Deadletter Queue Table",
            "helpText":"DLQ Table Ref: PROJECT:dataset.dlq",
            "paramType":"TEXT"
        },
        {
            "name":"autoscalingAlgorithm","label":"Autoscaling algorithm to use",
            "helpText":"Autoscaling algorithm to use: THROUGHPUT_BASED",
            "paramType":"TEXT",
            "isOptional":true
        },
        {
            "name":"numWorkers","label":"Number of workers Dataflow will start with",
            "helpText":"Number of workers Dataflow will start with",
            "paramType":"TEXT",
            "isOptional":true
        },

        {
            "name":"maxNumWorkers","label":"Maximum number of workers Dataflow job will use",
            "helpText":"Maximum number of workers Dataflow job will use",
            "paramType":"TEXT",
            "isOptional":true
        },
        {
            "name":"workerMachineType","label":"Worker Machine Type to use in Dataflow Job",
            "helpText":"Machine Type to Use: n1-standard-4",
            "paramType":"TEXT",
            "isOptional":true
        }
    ]},
    "sdk_info":{"language":"JAVA"}
}' > image_spec.json
gsutil cp image_spec.json ${TEMPLATE_IMAGE_SPEC}
rm image_spec.json
```

### Testing Template

The template unit tests can be run using:
```sh
mvn test
```

### Executing Template

The template requires the following parameters:
* outputDatasetTemplate: The name of the dataset or templated logic to extract it (ie. 'prefix_{schema_name}')
* outputTableNameTemplate: The name of the table or templated logic to extract it (ie. 'prefix_{table_name}')

The template has the following optional parameters:
* maxRetryAttempts: Max retry attempts, must be > 0. Default: no retries
* maxRetryDuration: Max retry duration in milliseconds, must be > 0. Default: no retries

Template can be executed using the following API call:
```sh
export JOB_NAME="${IMAGE_NAME}-`date +%Y%m%d-%H%M%S-%N`"
gcloud beta dataflow flex-template run ${JOB_NAME} \
        --project=${PROJECT} --region=us-central1 \
        --template-file-gcs-location=${TEMPLATE_IMAGE_SPEC} \
        --parameters inputSubscription=${SUBSCRIPTION},outputDeadletterTable=${DEADLETTER_TABLE}
```
