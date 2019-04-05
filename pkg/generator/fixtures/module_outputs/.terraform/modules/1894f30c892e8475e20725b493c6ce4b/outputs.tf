output "first_output" {
  description = "First output"
  value       = "${random_string.password.result}"
}
