data "aws-iam-policy-splitter_data_source" "this" {
  json_policies = [
    jsonencode({
      Version = "2012-10-17"
      Statement = [{
        Effect   = "Allow"
        Action   = ["s3:GetObject", "s3:PutObject"]
        Resource = "arn:aws:s3:::my-bucket/*"
      }]
    }),
    jsonencode({
      Version = "2012-10-17"
      Statement = [{
        Effect   = "Allow"
        Action   = ["ec2:DescribeInstances"]
        Resource = "*"
      }]
    }),
  ]

  bin_size = 6144
  max_bins = 10
}
