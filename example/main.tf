terraform {
  required_providers {
    ics = {
      source  = "UK2Group/ingenuitycloudservices"
      version = "~> 1.0"
    }
  }
}

provider "ics" {
  # Set your API token via environment variable: export ICS_API_TOKEN="your-token-here"
  # Alternatively, you can set it directly (not recommended for production):
  # api_token = "your-api-token-here"
}
