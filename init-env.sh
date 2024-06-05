#!/bin/bash

ENV_FILE=".env"

prompt_and_store_value() {
    read -p "Enter the $1: " value
    echo "$1:$value" >> $ENV_FILE
}

# Check if .env file already exists
if [ -f "$ENV_FILE" ]; then
    echo "Existing .env file found. If you proceed, it will be overwritten."
    read -p "Do you want to continue? (y/n): " overwrite
    if [ "$overwrite" != "y" ]; then
        echo "Aborted."
        exit 1
    fi
fi

echo "# Configuration for PR summary project" > $ENV_FILE

echo "GitHub Token: This token is needed to access GitHub API. You can generate one here: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token"
prompt_and_store_value "GITHUB_TOKEN"

echo "Gmail User: Enter your Gmail address"
prompt_and_store_value "GMAIL_USERNAME"

echo "Gmail Token: This token is needed to send emails via Gmail. You can follow this guide to obtain one: https://knowledge.workspace.google.com/kb/how-to-create-app-passwords-000009237#solution"
prompt_and_store_value "GMAIL_TOKEN"

echo "Target Repo: Enter the name of the target GitHub repository"
prompt_and_store_value "TARGET_REPOSITORY"

echo "Repo Owner: Enter the owner of the GitHub repository"
prompt_and_store_value "REPOSITORY_OWNER"


echo "Environment setup complete. Configuration saved to $ENV_FILE."

