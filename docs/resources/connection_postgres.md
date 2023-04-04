---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "materialize_connection_postgres Resource - terraform-provider-materialize"
subcategory: ""
description: |-
  The connection resource allows you to manage connections in Materialize.
---

# materialize_connection_postgres (Resource)

The connection resource allows you to manage connections in Materialize.

## Example Usage

```terraform
# Create a Postgres Connection
resource "materialize_connection_postgres" "example_postgres_connection" {
  name = "example_postgres_connection"
  host = "instance.foo000.us-west-1.rds.amazonaws.com"
  port = 5432
  user {
    secret {
      name          = "example"
      database_name = "database"
      schema_name   = "schema"
    }
  }
  password {
    name          = "example"
    database_name = "database"
    schema_name   = "schema"
  }
  database = "example"
}

# CREATE CONNECTION example_postgres_connection TO POSTGRES (
#     HOST 'instance.foo000.us-west-1.rds.amazonaws.com',
#     PORT 5432,
#     USER SECRET "database"."schema"."example"
#     PASSWORD SECRET "database"."schema"."example",
#     DATABASE 'example'
# );
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `database` (String) The target Postgres database.
- `host` (String) The Postgres database hostname.
- `name` (String) The identifier for the connection.
- `user` (Block List, Min: 1, Max: 1) The Postgres database username. (see [below for nested schema](#nestedblock--user))

### Optional

- `aws_privatelink` (Block List, Max: 1) The AWS PrivateLink configuration for the Postgres database. (see [below for nested schema](#nestedblock--aws_privatelink))
- `database_name` (String) The identifier for the connection database.
- `password` (Block List, Max: 1) The Postgres database password. (see [below for nested schema](#nestedblock--password))
- `port` (Number) The Postgres database port.
- `schema_name` (String) The identifier for the connection schema.
- `ssh_tunnel` (Block List, Max: 1) The SSH tunnel configuration for the Postgres database. (see [below for nested schema](#nestedblock--ssh_tunnel))
- `ssl_certificate` (Block List, Max: 1) The client certificate for the Postgres database. (see [below for nested schema](#nestedblock--ssl_certificate))
- `ssl_certificate_authority` (Block List, Max: 1) The CA certificate for the Postgres database. (see [below for nested schema](#nestedblock--ssl_certificate_authority))
- `ssl_key` (Block List, Max: 1) The client key for the Postgres database. (see [below for nested schema](#nestedblock--ssl_key))
- `ssl_mode` (String) The SSL mode for the Postgres database.

### Read-Only

- `connection_type` (String) The type of connection.
- `id` (String) The ID of this resource.
- `qualified_name` (String) The fully qualified name of the connection.

<a id="nestedblock--user"></a>
### Nested Schema for `user`

Optional:

- `secret` (Block List, Max: 1) The user secret value. (see [below for nested schema](#nestedblock--user--secret))
- `text` (String) The user text value.

<a id="nestedblock--user--secret"></a>
### Nested Schema for `user.secret`

Required:

- `name` (String) The user name.

Optional:

- `database_name` (String) The user database name.
- `schema_name` (String) The user schema name.



<a id="nestedblock--aws_privatelink"></a>
### Nested Schema for `aws_privatelink`

Required:

- `name` (String) The aws_privatelink name.

Optional:

- `database_name` (String) The aws_privatelink database name.
- `schema_name` (String) The aws_privatelink schema name.


<a id="nestedblock--password"></a>
### Nested Schema for `password`

Required:

- `name` (String) The password name.

Optional:

- `database_name` (String) The password database name.
- `schema_name` (String) The password schema name.


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
terraform import materialize_connection_postgres.example <connection_id>
```