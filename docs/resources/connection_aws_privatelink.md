---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "materialize_connection_aws_privatelink Resource - terraform-provider-materialize"
subcategory: ""
description: |-
  The connection resource allows you to manage connections in Materialize.
---

# materialize_connection_aws_privatelink (Resource)

The connection resource allows you to manage connections in Materialize.

## Example Usage

```terraform
# # Create a AWS Private Connection
# Note: you need the max_aws_privatelink_connections increased for this to work:
# show max_aws_privatelink_connections;
resource "materialize_connection_aws_privatelink" "example_privatelink_connection" {
  name               = "example_privatelink_connection"
  schema_name        = "public"
  service_name       = "com.amazonaws.us-east-1.materialize.example"
  availability_zones = ["use1-az2", "use1-az6"]
}

# CREATE CONNECTION example_privatelink_connection TO AWS PRIVATELINK (
#     SERVICE NAME 'com.amazonaws.us-east-1.materialize.example',
#     AVAILABILITY ZONES ('use1-az2', 'use1-az6')
# );
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `availability_zones` (List of String) The availability zones of the AWS PrivateLink service.
- `name` (String) The identifier for the connection.
- `service_name` (String) The name of the AWS PrivateLink service.

### Optional

- `database_name` (String) The identifier for the connection database.
- `schema_name` (String) The identifier for the connection schema.

### Read-Only

- `connection_type` (String) The type of connection.
- `id` (String) The ID of this resource.
- `qualified_name` (String) The fully qualified name of the connection.

## Import

Import is supported using the following syntax:

```shell
#Connections can be imported using the connection id:
terraform import materialize_connection_aws_privatelink.example <connection_id>
```