---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "materialize_connection_kafka Resource - terraform-provider-materialize"
subcategory: ""
description: |-
  The connection resource allows you to manage connections in Materialize.
---

# materialize_connection_kafka (Resource)

The connection resource allows you to manage connections in Materialize.

## Example Usage

```terraform
# Create a Kafka Connection
resource "materialize_connection_kafka" "example_kafka_connection" {
  name = "example_kafka_connection"
  kafka_broker {
    broker = "b-1.hostname-1:9096"
  }
  sasl_username = "example"
  sasl_password {
    name          = "kafka_password"
    database_name = "materialize"
    schema_name   = "public"
  }
  sasl_mechanisms = "SCRAM-SHA-256"
  progress_topic  = "example"
}

# CREATE CONNECTION database.schema.kafka_conn TO KAFKA (
#     BROKER 'example:9092'
#     PROGRESS TOPIC 'topic',
#     SASL MECHANISMS 'PLAIN',
#     SASL USERNAME 'user',
#     SASL PASSWORD SECRET "materialize"."public"."kafka_password"
# );

resource "materialize_connection_kafka" "example_kafka_connection_multiple_brokers" {
  name = "example_kafka_connection_multiple_brokers"
  kafka_broker {
    broker            = "b-1.hostname-1:9096"
    target_group_port = "9001"
    availability_zone = "use1-az1"
    privatelink_connection {
      name          = "example_aws_privatelink_conn"
      database_name = "materialize"
      schema_name   = "public"
    }
  }
  kafka_broker {
    broker            = "b-2.hostname-2:9096"
    target_group_port = "9002"
    availability_zone = "use1-az2"
    privatelink_connection {
      name          = "example_aws_privatelink_conn"
      database_name = "materialize"
      schema_name   = "public"
    }
  }
}

# CREATE CONNECTION materialize.public.example_kafka_connection_multiple_brokers TO KAFKA (
#     BROKERS (
#        'b-1.hostname-1:9096' USING AWS PRIVATELINK "materialize"."public"."example_aws_privatelink_conn" (PORT 9001, AVAILABILITY ZONE 'use1-az1'),
#        'b-2.hostname-2:9096' USING AWS PRIVATELINK "materialize"."public"."example_aws_privatelink_conn" (PORT 9002, AVAILABILITY ZONE 'use1-az2')
#     )
# );
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `kafka_broker` (Block List, Min: 1) The Kafka brokers configuration. (see [below for nested schema](#nestedblock--kafka_broker))
- `name` (String) The identifier for the connection.

### Optional

- `database_name` (String) The identifier for the connection database.
- `progress_topic` (String) The name of a topic that Kafka sinks can use to track internal consistency metadata.
- `sasl_mechanisms` (String) The SASL mechanism for the Kafka broker.
- `sasl_password` (Block List, Max: 1) The SASL password for the Kafka broker. (see [below for nested schema](#nestedblock--sasl_password))
- `sasl_username` (Block List, Max: 1) The SASL username for the Kafka broker. (see [below for nested schema](#nestedblock--sasl_username))
- `schema_name` (String) The identifier for the connection schema.
- `ssh_tunnel` (Block List, Max: 1) The SSH tunnel configuration for the Kafka broker. (see [below for nested schema](#nestedblock--ssh_tunnel))
- `ssl_certificate` (Block List, Max: 1) The client certificate for the Kafka broker. (see [below for nested schema](#nestedblock--ssl_certificate))
- `ssl_certificate_authority` (Block List, Max: 1) The CA certificate for the Kafka broker. (see [below for nested schema](#nestedblock--ssl_certificate_authority))
- `ssl_key` (Block List, Max: 1) The client key for the Kafka broker. (see [below for nested schema](#nestedblock--ssl_key))

### Read-Only

- `connection_type` (String) The type of connection.
- `id` (String) The ID of this resource.
- `qualified_name` (String) The fully qualified name of the connection.

<a id="nestedblock--kafka_broker"></a>
### Nested Schema for `kafka_broker`

Required:

- `broker` (String) The Kafka broker, in the form of `host:port`.

Optional:

- `availability_zone` (String) The availability zone of the Kafka broker.
- `privatelink_connection` (Block List, Max: 1) The AWS PrivateLink connection name in Materialize. (see [below for nested schema](#nestedblock--kafka_broker--privatelink_connection))
- `target_group_port` (Number) The port of the target group associated with the Kafka broker.

<a id="nestedblock--kafka_broker--privatelink_connection"></a>
### Nested Schema for `kafka_broker.privatelink_connection`

Required:

- `name` (String) The privatelink_connection name.

Optional:

- `database_name` (String) The privatelink_connection database name.
- `schema_name` (String) The privatelink_connection schema name.



<a id="nestedblock--sasl_password"></a>
### Nested Schema for `sasl_password`

Required:

- `name` (String) The sasl_password name.

Optional:

- `database_name` (String) The sasl_password database name.
- `schema_name` (String) The sasl_password schema name.


<a id="nestedblock--sasl_username"></a>
### Nested Schema for `sasl_username`

Optional:

- `secret` (Block List, Max: 1) The sasl_username secret value. (see [below for nested schema](#nestedblock--sasl_username--secret))
- `text` (String) The sasl_username text value.

<a id="nestedblock--sasl_username--secret"></a>
### Nested Schema for `sasl_username.secret`

Required:

- `name` (String) The sasl_username name.

Optional:

- `database_name` (String) The sasl_username database name.
- `schema_name` (String) The sasl_username schema name.



<a id="nestedblock--ssh_tunnel"></a>
### Nested Schema for `ssh_tunnel`

Required:

- `name` (String) The ssh_tunnel name.

Optional:

- `database_name` (String) The ssh_tunnel database name.
- `schema_name` (String) The ssh_tunnel schema name.


<a id="nestedblock--ssl_certificate"></a>
### Nested Schema for `ssl_certificate`

Optional:

- `secret` (Block List, Max: 1) The ssl_certificate secret value. (see [below for nested schema](#nestedblock--ssl_certificate--secret))
- `text` (String) The ssl_certificate text value.

<a id="nestedblock--ssl_certificate--secret"></a>
### Nested Schema for `ssl_certificate.secret`

Required:

- `name` (String) The ssl_certificate name.

Optional:

- `database_name` (String) The ssl_certificate database name.
- `schema_name` (String) The ssl_certificate schema name.



<a id="nestedblock--ssl_certificate_authority"></a>
### Nested Schema for `ssl_certificate_authority`

Optional:

- `secret` (Block List, Max: 1) The ssl_certificate_authority secret value. (see [below for nested schema](#nestedblock--ssl_certificate_authority--secret))
- `text` (String) The ssl_certificate_authority text value.

<a id="nestedblock--ssl_certificate_authority--secret"></a>
### Nested Schema for `ssl_certificate_authority.secret`

Required:

- `name` (String) The ssl_certificate_authority name.

Optional:

- `database_name` (String) The ssl_certificate_authority database name.
- `schema_name` (String) The ssl_certificate_authority schema name.



<a id="nestedblock--ssl_key"></a>
### Nested Schema for `ssl_key`

Required:

- `name` (String) The ssl_key name.

Optional:

- `database_name` (String) The ssl_key database name.
- `schema_name` (String) The ssl_key schema name.

## Import

Import is supported using the following syntax:

```shell
#Connections can be imported using the connection id:
terraform import materialize_connection_kafka.example <connection_id>
```