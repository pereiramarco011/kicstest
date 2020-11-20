resource "aws_config_configuration_aggregator" "account" {
  name = "example"

  account_aggregation_source {
    all_regions = true

  }
}

resource "aws_config_configuration_aggregator" "organization" {
  depends_on = [aws_iam_role_policy_attachment.organization]

  name = "example" # Required

  organization_aggregation_source {
    all_regions = true
    role_arn    = aws_iam_role.organization.arn
  }
}