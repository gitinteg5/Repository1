variable "port" {
    type = number
    default = 22
}

variable "cidr_block" {
    type = string
    default = "0.0.0.0/0"
}

variable "vpc_id" {
    type = string
    default = "some-vpc-id"
}