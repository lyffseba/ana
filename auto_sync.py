#!/usr/bin/env python3
"""
ANA Project Auto-Sync System
Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
"""

import os
import sys
import time
import json
import logging
import subprocess
from datetime import datetime
from pathlib import Path
from typing import Dict, List, Any, Optional

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('sync.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger('ana.sync')

class AutoSync:
    """Automated synchronization system"""
    
    def __init__(self):
        self.base_path = Path(os.getcwd())
        self.warp_session = "https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896"
        self.check_interval = 300  # 5 minutes
    
    def run_forever(self):
        """Run sync process continuously"""
        logger.info("Starting auto-sync system")
        
        while True:
            try:
                self.check_and_sync()
                time.sleep(self.check_interval)
            except KeyboardInterrupt:
                logger.info("Stopping auto-sync system")
                break
            except Exception as e:
                logger.error(f"Error in sync process: {e}")
                time.sleep(60)  # Wait a minute on error
    
    def check_and_sync(self):
        """Check for changes and sync"""
        logger.info("Checking for changes...")
        
        # Check git status
        if self._has_changes():
            logger.info("Changes detected")
            self._handle_changes()
        else:
            logger.info("No changes detected")
    
    def _has_changes(self) -> bool:
        """Check if there are any changes"""
        try:
            result = subprocess.run(
                ['git', 'status', '--porcelain'],
                capture_output=True,
                text=True
            )
            return bool(result.stdout.strip())
        except Exception as e:
            logger.error(f"Error checking git status: {e}")
            return False
    
    def _handle_changes(self):
        """Handle detected changes"""
        try:
            # Run tests
            self._run_tests()
            
            # Format code
            self._format_code()
            
            # Update documentation
            self._update_docs()
            
            # Commit and push
            self._commit_and_push()
            
        except Exception as e:
            logger.error(f"Error handling changes: {e}")
    
    def _run_tests(self):
        """Run project tests"""
        logger.info("Running tests...")
        
        try:
            subprocess.run(['go', 'test', './...'], check=True)
            logger.info("Tests passed")
        except subprocess.CalledProcessError as e:
            logger.error(f"Tests failed: {e}")
            raise
    
    def _format_code(self):
        """Format code"""
        logger.info("Formatting code...")
        
        try:
            subprocess.run(['go', 'fmt', './...'], check=True)
            logger.info("Code formatted")
        except subprocess.CalledProcessError as e:
            logger.error(f"Formatting failed: {e}")
            raise
    
    def _update_docs(self):
        """Update documentation"""
        logger.info("Updating documentation...")
        
        try:
            # Update all .md files with session reference
            for md_file in self.base_path.rglob('*.md'):
                self._update_doc_file(md_file)
            
            # Update Go files with session reference
            for go_file in self.base_path.rglob('*.go'):
                self._update_go_file(go_file)
            
            logger.info("Documentation updated")
        except Exception as e:
            logger.error(f"Error updating documentation: {e}")
            raise
    
    def _update_doc_file(self, file_path: Path):
        """Update a documentation file"""
        content = file_path.read_text()
        if self.warp_session not in content:
            new_content = f"# {file_path.stem}\nReference: {self.warp_session}\n\n{content}"
            file_path.write_text(new_content)
    
    def _update_go_file(self, file_path: Path):
        """Update a Go file"""
        content = file_path.read_text()
        if self.warp_session not in content:
            new_content = f"// {file_path.name}\n// Reference: {self.warp_session}\n\n{content}"
            file_path.write_text(new_content)
    
    def _commit_and_push(self):
        """Commit and push changes"""
        logger.info("Committing changes...")
        
        try:
            # Add all changes
            subprocess.run(['git', 'add', '.'], check=True)
            
            # Commit
            message = f"""chore: auto-sync update

- Run tests
- Format code
- Update documentation

Session: {self.warp_session}"""
            subprocess.run(['git', 'commit', '-m', message], check=True)
            
            # Push
            subprocess.run(['git', 'push'], check=True)
            
            logger.info("Changes pushed successfully")
        except subprocess.CalledProcessError as e:
            logger.error(f"Error in git operations: {e}")
            raise

def main():
    """Main function"""
    sync = AutoSync()
    sync.run_forever()

if __name__ == '__main__':
    main()
