output "codepipeline_name" {
  value = aws_codepipeline.pipeline.name
}

output "codepipeline_arn" {
  value = aws_codepipeline.pipeline.arn
}

output "s3_bucket_name" {
  value = aws_s3_bucket.codepipeline_artifacts.bucket
}

output "ec2_public_ip" {
  value = aws_instance.web_server.public_ip
}

output "github_connection_arn" {
  value = aws_codestarconnections_connection.github_connection.arn
}