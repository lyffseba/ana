#!/bin/bash
# Test script for repository
python -m pytest tests/ --cov=. --cov-report=xml
