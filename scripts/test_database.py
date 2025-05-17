#!/usr/bin/env python3
"""
Database connectivity test script for Ana project.
This script tests the PostgreSQL database connection and schema.
"""

import os
import sys
import psycopg2
from datetime import datetime, timedelta
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Database connection parameters from environment variables
DB_HOST = os.getenv("DB_HOST", "localhost")
DB_PORT = os.getenv("DB_PORT", "5432")
DB_NAME = os.getenv("DB_NAME", "ana_world")
DB_USER = os.getenv("DB_USER", "ana_user")
DB_PASSWORD = os.getenv("DB_PASSWORD", "your_secure_password")

def test_connection():
    """Test basic connection to the database."""
    print("Testing PostgreSQL connection...")
    
    try:
        conn = psycopg2.connect(
            host=DB_HOST,
            port=DB_PORT,
            dbname=DB_NAME,
            user=DB_USER,
            password=DB_PASSWORD
        )
        
        cursor = conn.cursor()
        cursor.execute("SELECT version();")
        db_version = cursor.fetchone()
        
        print(f"Successfully connected to PostgreSQL!")
        print(f"Server version: {db_version[0]}")
        
        cursor.close()
        return conn
    except Exception as e:
        print(f"Error connecting to PostgreSQL: {e}")
        sys.exit(1)

def test_tasks_table(conn):
    """Test if the tasks table exists and has the expected schema."""
    print("\nVerifying tasks table...")
    
    try:
        cursor = conn.cursor()
        
        # Check if tasks table exists
        cursor.execute("""
            SELECT EXISTS (
                SELECT FROM information_schema.tables 
                WHERE table_schema = 'public' 
                AND table_name = 'tasks'
            );
        """)
        
        table_exists = cursor.fetchone()[0]
        if not table_exists:
            print("Tasks table doesn't exist. Creating it...")
            create_tasks_table(cursor)
        else:
            print("Tasks table exists.")
            
            # Check columns to verify schema
            cursor.execute("""
                SELECT column_name, data_type 
                FROM information_schema.columns 
                WHERE table_schema = 'public' 
                AND table_name = 'tasks';
            """)
            
            columns = cursor.fetchall()
            print("\nTable schema:")
            for col in columns:
                print(f"  - {col[0]}: {col[1]}")
        
        cursor.close()
    except Exception as e:
        print(f"Error checking tasks table: {e}")
        sys.exit(1)

def create_tasks_table(cursor):
    """Create the tasks table with the expected schema."""
    try:
        cursor.execute("""
            CREATE TABLE tasks (
                id SERIAL PRIMARY KEY,
                title VARCHAR(255) NOT NULL,
                description TEXT,
                due_date TIMESTAMP,
                priority VARCHAR(50),
                project_id INTEGER,
                status VARCHAR(50),
                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            );
            
            CREATE INDEX idx_tasks_due_date ON tasks(due_date);
            CREATE INDEX idx_tasks_project_id ON tasks(project_id);
            CREATE INDEX idx_tasks_status ON tasks(status);
        """)
        
        cursor.connection.commit()
        print("Tasks table created successfully.")
    except Exception as e:
        print(f"Error creating tasks table: {e}")
        cursor.connection.rollback()
        sys.exit(1)

def insert_sample_data(conn):
    """Insert sample task data for testing."""
    print("\nInserting sample task data...")
    
    try:
        cursor = conn.cursor()
        
        # Check if we already have data
        cursor.execute("SELECT COUNT(*) FROM tasks;")
        count = cursor.fetchone()[0]
        
        if count > 0:
            print(f"Database already contains {count} tasks. Skipping sample data insertion.")
            return
        
        # Insert sample tasks
        tomorrow = datetime.now() + timedelta(days=1)
        next_week = datetime.now() + timedelta(days=7)
        
        cursor.execute("""
            INSERT INTO tasks (title, description, due_date, priority, project_id, status)
            VALUES 
            ('Reunión con Cliente', 'Discutir requisitos del proyecto para el nuevo edificio', %s, 'High', 1, 'To-Do'),
            ('Finalizar Planos', 'Completar la versión final de los planos', %s, 'Medium', 1, 'To-Do');
        """, (tomorrow, next_week))
        
        conn.commit()
        print("Sample data inserted successfully.")
        
        cursor.close()
    except Exception as e:
        print(f"Error inserting sample data: {e}")
        conn.rollback()

def main():
    """Main function to run database tests."""
    print("Ana PostgreSQL Database Test")
    print("===========================\n")
    
    conn = test_connection()
    test_tasks_table(conn)
    insert_sample_data(conn)
    
    print("\nAll database tests completed successfully!")
    conn.close()

if __name__ == "__main__":
    main()

