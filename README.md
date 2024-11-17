# Enabling GitHub Actions Bot to Modify Your Repository

## Motivation

When I started using [golangci-lint](https://github.com/golangci/golangci-lint), I found my code violates [WSL](https://github.com/bombsimon/wsl) suggestions among others of course.

I learned from [here](https://github.com/bombsimon/wsl/tree/c862f085c18f8560c5aa50183cb4fbb9a11656c3?tab=readme-ov-file#usage) that I can't use WSL's `--fix` functionality through golangci-lint. According to WSL's documentation, this means if I want automatic fixes for WSL's suggestions, I have to install WSL as a standalone tool.

So maybe its best to install all the golangci-lint linter's myself, perhaps they all have one-off quirks like this that prevents golangci-lint to cleanup my code.

Maybe installing/maintaining them on github image using CI is better than maintainining them on my laptop.

Ok, so If i want all the linters on remote github runner then I should allow the runner to update my code with all the linters's fixes.

So, here I fumbled my way through getting github actions to commit to my repo and its simple in hindsight (create a personal access token and provide it as a secret in the repo), but arriving to this conclusion was time consuing.

## TL;DR

1. Create a token under personal account from https://github.com/settings/tokens, name it `WORKFLOW_TOKEN` (for example). Assume the token's value is `ghp_npXyyg50n1taucGKq1YGLvrEChfzz71Hnl32`
1. Create a GitHub org and add your personal account to it
1. Create org level secret with name `WORKFLOW_TOKEN` and value `ghp_npXyyg50n1taucGKq1YGLvrEChfzz71Hnl32`
1. Create repository
1. Add `.github/workflows/ci.yml`
1. Test with: `git commit --allow-empty -m "trigger build" && git push`
1. Verify with `git pull` to see `github-actions[bot]`'s commit

## Why Use Organization Secrets instead of personal account secrets?

Using organizational secrets instead of individual repository secrets offers several advantages:

1. **Centralized Secret Management**: Organization secrets can be shared across all repositories within the organization, eliminating the need to manually copy and paste the `WORKFLOW_TOKEN` for each new repository.

1. **Consistency**: All repositories in the organization can use the same token, ensuring consistent access levels and reducing configuration errors.
1. **Maintenance**: When you need to rotate or update the token, you only need to update it in one place (the organization settings) rather than updating it in each repository individually.
1. **Access Control**: Organization administrators can control which repositories have access to the secrets, providing better security management.

## Detailed Guide

### 1. Create Personal Access Token

Navigate to the token creation page:

https://github.com/settings, then:

```
→ Developer settings
→ Personal access tokens
→ Tokens (classic)
→ Generate new token (classic)
```

Direct URL: https://github.com/settings/tokens

1. Create a new token named `WORKFLOW_TOKEN`
   (Your token will have format similar to `ghp_npXyyg50n1taucGKq1YGLvrEChfzz71Hnl32`)
2. Required `WORKFLOW_TOKEN` permissions: TBD (to be documented in next update)

### 2. Configure Organization Secret

Navigate to organization secrets:

```
Your profile
→ Your organizations
→ [Organization name] link
→ Settings
→ Secrets and variables
→ Actions

Direct URL: https://github.com/organizations/[org-name]/settings/secrets/actions
```

Create new organizational secret:

- Name: `WORKFLOW_TOKEN`
- Value: (paste your token that has format similar to `ghp_npXyyg50n1taucGKq1YGLvrEChfzz71Hnl32`)

### 3. Configure GitHub Actions

Create `.github/workflows/ci.yml` with the following content. This example uses [Federico Grandi's](https://github.com/EndBug) [Add & Commit](https://github.com/marketplace/actions/add-commit) Github action to handle the commit process.

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

Note that we include `[skip ci]` in the commit message to prevent an infinite loop of workflow runs.

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

- This generates conflicts often when using renovate too since they're both updating same files at once, nee to search for workaround. To reporduce, change this `EndBug/add-and-commit@a94899bca583c204427a224a7af87c02f9b325d5 # v9` to this `EndBug/add-and-commit@v9` for example and observe conflict in pull request.
- Document specific token permissions needed
- Add troubleshooting section
- Add security considerations
- Add best practices for workflow file
