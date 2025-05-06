<!--# go-mercado-pago-ms - Mercado Pago Integration Microservice

[![Go Version](https://img.shields.io/github/go-mod/go-version/felipemacedo1/go-mercado-pago-ms)](https://golang.org/)
[![Build Status](https://github.com/felipemacedo1/go-mercado-pago-ms/actions/workflows/go.yml/badge.svg)](https://github.com/felipemacedo1/go-mercado-pago-ms/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/felipemacedo1/go-mercado-pago-ms/badge.svg?branch=main)](https://coveralls.io/github/felipemacedo1/go-mercado-pago-ms?branch=main)

This project is a microservice written in Go that integrates with the Mercado Pago API to handle payment processing. It's designed to be clean, maintainable, and follows best practices for microservice development in Go.

This microservice was built following the guidance and development support of the Google Gemini AI model.

## 🚀 Features

*   Create payment preferences via the `/checkout` endpoint.
*   Receive and process Mercado Pago webhook notifications on the `/webhook` endpoint.
*   Retrieve payment status using the `/payment/:id` endpoint.
*   Environment variable configuration using `.env`.
*   Custom structured logging with Zap.
*   Mock JWT authentication middleware.
*   Containerization with Docker.

## 🛠️ Technologies Used
*   **Go:** The primary programming language.
*   **Echo:** A high-performance, minimalist Go web framework.
*   **Mercado Pago API:** For payment processing.
*   **godotenv:** To load environment variables from a `.env` file.
*   **Zap:** A fast, structured, leveled logging library.
*   **httptest:** Go's built-in package for HTTP testing.
*   **testify:** A Go testing toolkit.
*   **Docker:** For containerization.





-->