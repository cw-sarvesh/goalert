# High Priority Alerts

GoAlert can elevate specific alerts to voice calls based on alert metadata.

## Configuration

Set the label key and value used to mark an alert as high priority:

```json
{
  "Alerts": {
    "HighPriorityLabelKey": "alerts/priority",
    "HighPriorityLabelValue": "high"
  }
}
```

## Usage

When an alert includes matching metadata, GoAlert routes the notification to the user's voice contact method, bypassing Slack and other contact methods.

Example GraphQL mutation creating a high-priority alert:

```graphql
mutation {
  createAlert(input:{
    serviceID:"<service-id>",
    summary:"database down",
    meta:[{key:"alerts/priority", value:"high"}]
  }) { id }
}
```
