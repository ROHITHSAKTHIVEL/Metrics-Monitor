#!/bin/bash
set -e

echo "Checking if database exists..."
psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'metrics_db'" | grep -q 1 || psql -U postgres -c "CREATE DATABASE metrics_db"

echo "Database check completed."
