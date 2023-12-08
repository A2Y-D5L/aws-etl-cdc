# aws-etl-cdc

A custom Extract-Transform-Load (ETL) and Change Data Capture (CDC) system built with serverless AWS services. It is designed as a proof-of-concept for use cases where DynamoDB streams aren't an option.

## Requirements

### Trigger

The process must be triggered by an SNS notification about an Item in a DynamoDB Table being created, updated or deleted.

### Notification Durability

The system must prevent against notifications being lost for any reason.

### Order

SNS notifications must be processed in the order they arrive.

### Extract-Transform-Load (ETL)

- Data must be extracted from multiple data stores (DynamoDB and S3).
- System must be able to handle extracted data sets that exceed the SQS message size limit.
- Extracted data must be transformed into multiple unique data sets according to destination-specific data models.
- Transformations must be performed concurrently.
- Transformed data must be loaded into temporary persistence.
- Loads into temporary persistence must be performed concurrently.

### Transmission

- Successful loads into temporary persistence must trigger transmission of the loaded data.
- Transmission of destination-specific data sets must be transmitted to a specified destination.
- Transmission destination must be a resource-oriented HTTP API.
- Payloads must be sent as JSON and conform to the destination API's OpenAPI spec.
- Transmissions must be performed concurrently and independent of one another.

### Transactions

- All transmissions associated with the same triggering notification must be treated as a single transaction.
- If a failure occurs for one of several transmissions derived from the same notifcation, then any previously successful transmissions derived from the same notification must be rolled back.
- Transmission failures must be persisted along with the triggering notification and all associated transformations.
