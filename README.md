# autoweb

**autoweb** is a lightweight Go-based container that serves static HTML content from a Git repository and automatically updates the content when a webhook triggers a `git pull`. This container is designed to be deployed in Kubernetes or other container environments where you want to host static websites and sync them with the latest updates from a remote Git repository.

## Features

* Serves static HTML from the `public/` directory of a cloned repository.
* Automatically pulls updates from the Git repository when triggered by a webhook.
* Minimal footprint with Go-based implementation for both Git operations and serving static content.
* Easy to integrate with Continuous Deployment setups using webhooks.


## Environment Variables

`REPO_URL:` The URL of the Git repository to clone. This repository should have the `public/` directory containing the pre-generated static content.

## How it Works

1. **Clone the Repo**: On startup, the container clones the Git repository defined by the `REPO_URL` environment variable if it hasn't already.
2. **Serve Static Content**: The container serves the static content found in the `public/` directory.
3. **Webhook Integration**: The container listens for POST requests on `/webhook`. When a request is received, the container pulls the latest changes from the Git repository and updates the static content being served.
