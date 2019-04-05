# terraform-variables-generator

Simple tool to generate variables file or module outputs from existing terraform configuration.

## Installation

```bash
go get -u github.com/alexandrst88/terraform-variables-generator
```

## Usage

```
terraform-variables-generator --help
```

### Variables

Variable generation is default:

```bash
terraform-variables-generator
```

Set name of generated file:
```bash
terraform-variables-generator --vars-file=./some_name.tf
```

It will find all `*.tf` files in current directory, and generate a file. If you already have this file, it will ask to override it.

### Modules outputs

Generate only modules outputs:

```bash
terraform-variables-generator  --vars=false --module-outputs=true
```

Generate outputs only for some modules:
```bash
terraform-variables-generator  --vars=false --module-outputs --modules-filter "^module-name$"
```

Module output generation will find all outputs for matching modules and re-output them in root terraform config with output names prefixed by module name.

### Example

```hcl
resource "aws_vpc" "vpc" {
  cidr_block           = "${var.cidr}"
  enable_dns_hostnames = "${var.enable_dns_hostnames}"
  enable_dns_support   = "${var.enable_dns_support}"

  tags {
    Name = "${var.name}"
  }
}
```

Will generate:

```hcl
variable "cidr" {
  description = ""
}

variable "enable_dns_hostnames" {
  description = ""
}

variable "enable_dns_support" {
  description = ""
}

variable "name" {
  description = ""
}
```

Module:
```hcl
module "prod" {
  source = ./
}
```

With outputs like this:
```hcl
output "vpc_id" {
  value = "${aws_vpc.main.id}"
}
```

Will generate outputs like this:
```hcl
output "prod_vpc_id" {
  description = ""
  value       = "${module.prod.vpc_id}"
}
```

## Tests

Run tests and linter

```bash
go vet ./...
go test -v -race ./...
golint -set_exit_status $(go list ./...)
golangci-lint run
```
