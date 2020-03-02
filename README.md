# Bitbucket Cascade Merge

Bitbucket Cascade Merge is a service used to automatically cascade changes
after a pull request has been successfully merged. This feature has not
been ported to Bitbucket Cloud.

You can show you interest and vote for this feature :
[BCLOUD-14286](https://jira.atlassian.com/browse/BCLOUD-14286)

## Usage

### What do you need ?

* A service to host this service (eg. Amazon ECS)
* An account on bitbucket.org with r/w privileges to the project

### Get an Application Password

The application password will be used to call the Bitbucket API and fetch/push
the repository.

1. Open https://bitbucket.org/account/user/{{username}}
2. Go to Access Management > App passwords
3. Click on **Create an app password**
4. Type a **label** and select the following permissions :
   * Repositories : Read, Write
   * Pull Requests : Read, Write
5. Copy the password somewhere safe, you will need it later

### Configure a webhook on the repository

1. Navigate to the repository
2. Go to Settings > Workflow > Webhooks
3. Click on **Add webhook**
4. Type a title, the url of your container and select
   **Choose from a full list of triggers** : Pull Request > Merged

### Configure the container

The container can be configured with environment variable.

Key | Default Value | Description
--- | --- | ---
PORT | 5000 | Server will listen on this port
USERNAME | | Bitbucket username
PASSWORD | | Bitbucket app password

## Requirements

[Libgit2 v0.25](https://github.com/libgit2/libgit2/archive/v0.25.0.tar.gz)
must be installed if you do not use the Docker image provided.

