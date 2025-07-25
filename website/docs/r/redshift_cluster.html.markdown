---
subcategory: "Redshift"
layout: "aws"
page_title: "AWS: aws_redshift_cluster"
description: |-
  Provides a Redshift Cluster resource.
---

# Resource: aws_redshift_cluster

Provides a Redshift Cluster Resource.

~> **Note:** All arguments including the username and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](https://www.terraform.io/docs/state/sensitive-data.html).

~> **NOTE:** A Redshift cluster's default IAM role can be managed both by this resource's `default_iam_role_arn` argument and the [`aws_redshift_cluster_iam_roles`](redshift_cluster_iam_roles.html) resource's `default_iam_role_arn` argument. Do not configure different values for both arguments. Doing so will cause a conflict of default IAM roles.

-> **Note:** Write-Only argument `master_password_wo` is available to use in place of `master_password`. Write-Only arguments are supported in HashiCorp Terraform 1.11.0 and later. [Learn more](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments).

## Example Usage

### Basic Usage

```terraform
resource "aws_redshift_cluster" "example" {
  cluster_identifier = "tf-redshift-cluster"
  database_name      = "mydb"
  master_username    = "exampleuser"
  master_password    = "Mustbe8characters"
  node_type          = "dc1.large"
  cluster_type       = "single-node"
}
```

### With Managed Credentials

```terraform
resource "aws_redshift_cluster" "example" {
  cluster_identifier = "tf-redshift-cluster"
  database_name      = "mydb"
  master_username    = "exampleuser"
  node_type          = "dc1.large"
  cluster_type       = "single-node"

  manage_master_password = true
}
```

## Argument Reference

This resource supports the following arguments:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `cluster_identifier` - (Required) The Cluster Identifier. Must be a lower case string.
* `database_name` - (Optional) The name of the first database to be created when the cluster is created.
  If you do not provide a name, Amazon Redshift will create a default database called `dev`.
* `default_iam_role_arn` - (Optional) The Amazon Resource Name (ARN) for the IAM role that was set as default for the cluster when the cluster was created.
* `node_type` - (Required) The node type to be provisioned for the cluster.
* `cluster_type` - (Optional) The cluster type to use. Either `single-node` or `multi-node`.
* `manage_master_password` - (Optional) Whether to use AWS SecretsManager to manage the cluster admin credentials.
  Conflicts with `master_password` and `master_password_wo`.
  One of `master_password` or `manage_master_password` is required unless `snapshot_identifier` is provided.
* `master_password` - (Optional) Password for the master DB user.
  Conflicts with `manage_master_password` and `master_password_wo`.
  One of `master_password`, `master_password_wo` or `manage_master_password` is required unless `snapshot_identifier` is provided.
  Note that this may show up in logs, and it will be stored in the state file.
  Password must contain at least 8 characters and contain at least one uppercase letter, one lowercase letter, and one number.
* `master_password_wo` - (Optional, Write-Only) Password for the master DB user.
  Conflicts with `manage_master_password` and `master_password`.
  One of `master_password_wo`, `master_password` or `manage_master_password` is required unless `snapshot_identifier` is provided.
  Note that this may show up in logs.
  Password must contain at least 8 characters and contain at least one uppercase letter, one lowercase letter, and one number.
* `master_password_wo_version` - (Optional) Used together with `master_password_wo` to trigger an update. Increment this value when an update to the `master_password_wo` is required.
* `master_password_secret_kms_key_id` - (Optional) ID of the KMS key used to encrypt the cluster admin credentials secret.
* `master_username` - (Required unless a `snapshot_identifier` is provided) Username for the master DB user.
* `multi_az` - (Optional) Specifies if the Redshift cluster is multi-AZ.
* `vpc_security_group_ids` - (Optional) A list of Virtual Private Cloud (VPC) security groups to be associated with the cluster.
* `cluster_subnet_group_name` - (Optional) The name of a cluster subnet group to be associated with this cluster. If this parameter is not provided the resulting cluster will be deployed outside virtual private cloud (VPC).
* `availability_zone` - (Optional) The EC2 Availability Zone (AZ) in which you want Amazon Redshift to provision the cluster. For example, if you have several EC2 instances running in a specific Availability Zone, then you might want the cluster to be provisioned in the same zone in order to decrease network latency. Can only be changed if `availability_zone_relocation_enabled` is `true`.
* `availability_zone_relocation_enabled` - (Optional) If true, the cluster can be relocated to another availabity zone, either automatically by AWS or when requested. Default is `false`. Available for use on clusters from the RA3 instance family.
* `preferred_maintenance_window` - (Optional) The weekly time range (in UTC) during which automated cluster maintenance can occur.
  Format: ddd:hh24:mi-ddd:hh24:mi
* `cluster_parameter_group_name` - (Optional) The name of the parameter group to be associated with this cluster.
* `automated_snapshot_retention_period` - (Optional) The number of days that automated snapshots are retained. If the value is 0, automated snapshots are disabled. Even if automated snapshots are disabled, you can still create manual snapshots when you want with create-cluster-snapshot. Default is 1.
* `port` - (Optional) The port number on which the cluster accepts incoming connections. Valid values are between `1115` and `65535`.
  The cluster is accessible only via the JDBC and ODBC connection strings.
  Part of the connection string requires the port on which the cluster will listen for incoming connections.
  Default port is `5439`.
* `cluster_version` - (Optional) The version of the Amazon Redshift engine software that you want to deploy on the cluster.
  The version selected runs on all the nodes in the cluster.
* `allow_version_upgrade` - (Optional) If true , major version upgrades can be applied during the maintenance window to the Amazon Redshift engine that is running on the cluster. Default is `true`.
* `apply_immediately` - (Optional) Specifies whether any cluster modifications are applied immediately, or during the next maintenance window. Default is `false`.
* `aqua_configuration_status` - (Optional, **Deprecated**) The value represents how the cluster is configured to use AQUA (Advanced Query Accelerator) after the cluster is restored.
  No longer supported by the AWS API.
  Always returns `auto`.
* `number_of_nodes` - (Optional) The number of compute nodes in the cluster. This parameter is required when the ClusterType parameter is specified as multi-node. Default is 1.
* `publicly_accessible` - (Optional) If true, the cluster can be accessed from a public network. Default is `false`.
* `encrypted` - (Optional) If true , the data in the cluster is encrypted at rest.
  Default is `true`.
* `enhanced_vpc_routing` - (Optional) If true , enhanced VPC routing is enabled.
* `kms_key_id` - (Optional) The ARN for the KMS encryption key. When specifying `kms_key_id`, `encrypted` needs to be set to true.
* `elastic_ip` - (Optional) The Elastic IP (EIP) address for the cluster.
* `skip_final_snapshot` - (Optional) Determines whether a final snapshot of the cluster is created before Amazon Redshift deletes the cluster. If true , a final cluster snapshot is not created. If false , a final cluster snapshot is created before the cluster is deleted. Default is false.
* `final_snapshot_identifier` - (Optional) The identifier of the final snapshot that is to be created immediately before deleting the cluster. If this parameter is provided, `skip_final_snapshot` must be false.
* `snapshot_arn` - (Optional) The ARN of the snapshot from which to create the new cluster. Conflicts with `snapshot_identifier`.
* `snapshot_identifier` - (Optional) The name of the snapshot from which to create the new cluster.  Conflicts with `snapshot_arn`.
* `snapshot_cluster_identifier` - (Optional) The name of the cluster the source snapshot was created from.
* `owner_account` - (Optional) The AWS customer account used to create or copy the snapshot. Required if you are restoring a snapshot you do not own, optional if you own the snapshot.
* `iam_roles` - (Optional) A list of IAM Role ARNs to associate with the cluster. A Maximum of 10 can be associated to the cluster at any time.
* `maintenance_track_name` - (Optional) The name of the maintenance track for the restored cluster. When you take a snapshot, the snapshot inherits the MaintenanceTrack value from the cluster. The snapshot might be on a different track than the cluster that was the source for the snapshot. For example, suppose that you take a snapshot of  a cluster that is on the current track and then change the cluster to be on the trailing track. In this case, the snapshot and the source cluster are on different tracks. Default value is `current`.
* `manual_snapshot_retention_period` - (Optional)  The default number of days to retain a manual snapshot. If the value is -1, the snapshot is retained indefinitely. This setting doesn't change the retention period of existing snapshots. Valid values are between `-1` and `3653`. Default value is `-1`.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

For more detailed documentation about each argument, refer to
the [AWS official documentation](http://docs.aws.amazon.com/cli/latest/reference/redshift/index.html#cli-aws-redshift).

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - Amazon Resource Name (ARN) of cluster
* `id` - The Redshift Cluster ID.
* `cluster_identifier` - The Cluster Identifier
* `cluster_type` - The cluster type
* `node_type` - The type of nodes in the cluster
* `database_name` - The name of the default database in the Cluster
* `availability_zone` - The availability zone of the Cluster
* `automated_snapshot_retention_period` - The backup retention period
* `preferred_maintenance_window` - The backup window
* `endpoint` - The connection endpoint
* `encrypted` - Whether the data in the cluster is encrypted
* `vpc_security_group_ids` - The VPC security group Ids associated with the cluster
* `dns_name` - The DNS name of the cluster
* `master_password_secret_arn` - ARN of the cluster admin credentials secret
* `port` - The Port the cluster responds on
* `cluster_version` - The version of Redshift engine software
* `cluster_parameter_group_name` - The name of the parameter group to be associated with this cluster
* `cluster_subnet_group_name` - The name of a cluster subnet group to be associated with this cluster
* `cluster_public_key` - The public key for the cluster
* `cluster_revision_number` - The specific revision number of the database in the cluster
* `cluster_nodes` - The nodes in the cluster. Cluster node blocks are documented below
* `cluster_namespace_arn` - The namespace Amazon Resource Name (ARN) of the cluster
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

Cluster nodes (for `cluster_nodes`) support the following attributes:

* `node_role` - Whether the node is a leader node or a compute node
* `private_ip_address` - The private IP address of a node within a cluster
* `public_ip_address` - The public IP address of a node within a cluster

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

- `create` - (Default `75m`)
- `update` - (Default `75m`)
- `delete` - (Default `40m`)

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Redshift Clusters using the `cluster_identifier`. For example:

```terraform
import {
  to = aws_redshift_cluster.myprodcluster
  id = "tf-redshift-cluster-12345"
}
```

Using `terraform import`, import Redshift Clusters using the `cluster_identifier`. For example:

```console
% terraform import aws_redshift_cluster.myprodcluster tf-redshift-cluster-12345
```
