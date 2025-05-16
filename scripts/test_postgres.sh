#!/bin/bash
# PostgreSQL database setup and test script for Ana project

# Set text colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Ana PostgreSQL Database Setup & Test${NC}"
echo -e "${YELLOW}=================================${NC}\n"

# Load environment variables from .env file
if [ -f ../.env ]; then
  export $(grep -v '^#' ../.env | xargs)
fi

DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"5432"}
DB_NAME=${DB_NAME:-"ana_world"}
DB_USER=${DB_USER:-"ana_user"}
DB_PASSWORD=${DB_PASSWORD:-"your_secure_password"}

# Function to check if PostgreSQL is installed
check_postgres() {
  if ! command -v psql &> /dev/null; then
    echo -e "${RED}PostgreSQL client not found. Please install PostgreSQL.${NC}"
    echo "On Ubuntu/Debian: sudo apt install postgresql postgresql-contrib"
    echo "On macOS with Homebrew: brew install postgresql"
    exit 1
  fi
  
  echo -e "${GREEN}PostgreSQL client found.${NC}"
  
  # Check if PostgreSQL server is running
  if ! pg_isready -h $DB_HOST -p $DB_PORT &> /dev/null; then
    echo -e "${RED}PostgreSQL server is not running on $DB_HOST:$DB_PORT.${NC}"
    echo "Please start the PostgreSQL service."
    echo "On Ubuntu/Debian: sudo systemctl start postgresql"
    echo "On macOS with Homebrew: brew services start postgresql"
    exit 1
  fi
  
  echo -e "${GREEN}PostgreSQL server is running.${NC}"
}

# Function to create database and user
setup_database() {
  echo -e "\n${YELLOW}Setting up database...${NC}"
  
  # Check if the user can connect to PostgreSQL as postgres
  if ! psql -h $DB_HOST -p $DB_PORT -U postgres -c '\conninfo' &> /dev/null; then
    echo -e "${RED}Cannot connect to PostgreSQL as postgres user.${NC}"
    echo "You need to run the following commands manually as a PostgreSQL superuser:"
    echo "--------------------------------------------------------------"
    echo "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';"
    echo "ALTER ROLE $DB_USER SET client_encoding TO 'utf8';"
    echo "ALTER ROLE $DB_USER SET default_transaction_isolation TO 'read committed';"
    echo "ALTER ROLE $DB_USER SET timezone TO 'UTC';"
    echo "CREATE DATABASE $DB_NAME OWNER $DB_USER;"
    echo "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"
    echo "--------------------------------------------------------------"
    
    read -p "Have you run these commands? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      echo "Exiting. Please run the commands and try again."
      exit 1
    fi
  else
    # Create user and database
    echo "Creating user $DB_USER..."
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null || echo "User may already exist"
    
    echo "Setting user properties..."
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "ALTER ROLE $DB_USER SET client_encoding TO 'utf8';" 2>/dev/null || true
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "ALTER ROLE $DB_USER SET default_transaction_isolation TO 'read committed';" 2>/dev/null || true
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "ALTER ROLE $DB_USER SET timezone TO 'UTC';" 2>/dev/null || true
    
    echo "Creating database $DB_NAME..."
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;" 2>/dev/null || echo "Database may already exist"
    
    echo "Granting privileges..."
    psql -h $DB_HOST -p $DB_PORT -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;" 2>/dev/null || true
  fi
  
  # Check if we can connect with the new user
  if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\conninfo' &> /dev/null; then
    echo -e "${GREEN}Successfully connected with $DB_USER user.${NC}"
  else
    echo -e "${RED}Could not connect with $DB_USER user. Check your credentials.${NC}"
    exit 1
  fi
}

# Function to run migrations
run_migrations() {
  echo -e "\n${YELLOW}Running migrations...${NC}"
  
  # Check if migration file exists
  migration_file="../db/migrations/0001_init_tasks.sql"
  if [ ! -f "$migration_file" ]; then
    echo -e "${RED}Migration file not found: $migration_file${NC}"
    exit 1
  fi
  
  # Run migration
  echo "

