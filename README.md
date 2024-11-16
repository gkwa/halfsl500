# Enabling GitHub Actions Bot to Modify Your Repository

## TL;DR

1. Create a Github organization and make your account a member of it
1. Create personal access token at https://github.com/settings/tokens
1. Create organizational secret `WORKFLOW_TOKEN` with the token value
1. Create repository
1. Add `.github/workflows/ci.yml` (example below)
1. Test with: `git commit --allow-empty -m "trigger build" && git push`
1. Verify with `git pull` to see `github-actions[bot]`'s commit

## Why Use Organization Secrets instead of personal account secrets?

Using organizational secrets instead of individual repository secrets offers several advantages:

1. **Centralized Secret Management**: Organization secrets can be shared across all repositories within the organization, eliminating the need to manually copy and paste the `WORKFLOW_TOKEN` for each new repository.

2. **Consistency**: All repositories in the organization can use the same token, ensuring consistent access levels and reducing configuration errors.

3. **Maintenance**: When you need to rotate or update the token, you only need to update it in one place (the organization settings) rather than updating it in each repository individually.

4. **Access Control**: Organization administrators can control which repositories have access to the secrets, providing better security management.

## Detailed Guide

### 1. Create Personal Access Token

1. Navigate to https://github.com/settings/tokens using these steps:
   ```
   1. Login to GitHub
   2. Click on your profile picture
   3. Click the Settings link
   4. Click "Developer settings" link (bottom of left sidebar)
   5. Click "Personal access tokens" dropdown
   6. Click "Tokens (classic)" item
   7. Click "Generate new token" combo box
   8. Click "Generate new token (classic)" item
   ```
2. Create a new token named `WORKFLOW_TOKEN`
   (Your token will have format similar to this: `ghp_npXyyg50n1taucGKq1YGLvrEChfzz71Hnl32`)
3. Required `WORKFLOW_TOKEN` permissions: TBD (to be documented in next update)

### 2. Configure Organization Secret

1. Navigate to your organization settings using these steps:

   ```
   1. Login to GitHub
   2. Click on your profile picture
   3. Click "Your organizations" link
   4. Click on your organization name link
   5. Click Settings link
      (Example URL: https://github.com/organizations/gkwa/settings/profile)
   6. Click "Secrets and variables" link
   7. Click "Actions" link
      (Example URL: https://github.com/organizations/gkwa/settings/secrets/actions)
   8. Click "New organization secret" button
   ```

   Example organization settings URL: https://github.com/organizations/gkwa/settings/profile (replace 'gkwa' with your organization name)

   Example Actions settings URL: https://github.com/organizations/gkwa/settings/secrets/actions (replace 'gkwa' with your organization name)

2. Create new organizational secret:
   - Name: `WORKFLOW_TOKEN`
   - Value: (paste your token that has format similar to this: `ghp_npXyyg50n1taucGKq1YGLvrEChfzz71Hnl32`)

### 3. Configure GitHub Actions

Create `.github/workflows/ci.yml` with the following content. This example uses the Github Action [EndBug/add-and-commit](https://github.com/EndBug/add-and-commit?tab=readme-ov-file#add--commit) to handle the commit process.

Note that we include `[skip ci]` in the commit message to prevent an infinite loop of workflow runs:

```yaml
name: Build & Test
on:
  push:
    branches:
      - "*"
  pull_request:
    branches:
      - "*"
jobs:
  test:
    name: Build & Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@cbb722410c2e876e24abbe8de2cc27693e501dcb
        with:
          token: ${{ secrets.WORKFLOW_TOKEN }}
          ref: ${{ github.head_ref || github.ref_name }}
      - run: |
          date >>data.txt # arbitrary change to ensure our change is pushed back to repo
      - run: |
          cat data.txt # placeholder for validation check
      - uses: EndBug/add-and-commit@a94899bca583c204427a224a7af87c02f9b325d5 # v9
        with:
          add: .
          message: "chore: append date to data.txt [skip ci]" # [skip ci] prevents infinite workflow runs
          default_author: github_actions
          github_token: ${{ secrets.WORKFLOW_TOKEN }}
```

Uses `EndBug`'s `add-and-commit` github action:
https://github.com/EndBug/add-and-commit?tab=readme-ov-file#add--commit

### 4. Testing the Setup

1. Trigger the workflow:

```bash
git commit --allow-empty -m "trigger build" && git push
```

2. Pull the changes to see the bot's commit:

```bash
git pull
```

3. Example of a successful bot commit (from [example repository](https://github.com/gkwa/halfsl500/commit/3fb3a662d0ecd77d0615234f5a3dfc3e3cfb7411)):

```
commit 3fb3a662d0ecd77d0615234f5a3dfc3e3cfb7411
Author: github-actions <41898282+github-actions[bot]@users.noreply.github.com>
Date:   Fri Nov 15 23:00:17 2024 +0000
    chore: append date to data.txt [skip ci]
diff --git a/data.txt b/data.txt
index 1e15dd7..cf49456 100644
--- a/data.txt
@@ -10,3 +10,4 @@ Fri Nov 15 21:10:26 UTC 2024
 Fri Nov 15 21:13:39 UTC 2024
 Fri Nov 15 21:15:00 UTC 2024
 Fri Nov 15 21:19:51 UTC 2024
+Fri Nov 15 23:00:15 UTC 2024
```

You can see more examples of successful bot commits in the [example repository](https://github.com/gkwa/halfsl500/commits/master/).

### Todo

- Document specific token permissions needed
- Add troubleshooting section
- Add security considerations
- Add best practices for workflow file
