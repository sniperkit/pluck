hub is a command-line wrapper for git that makes you better at GitHub.

# Install HUB

## macOS: install the latest release
```bash
$ brew install hub
```

## Other platforms: fetch a precompiled binary release, or build your own from source.

You will need Go development environment.
```bash
$ git clone https://github.com/github/hub.git && cd hub
```

## assuming that `~/bin` is in your PATH:
```bash
$ script/build -o ~/bin/hub
```

## alias it as git
```bash
$ alias git=hub
$ git version
git version 2.17.0
hub version 2.3.0 # ← it works!
```

## As a contributor to open-source
Whether you are beginner or an experienced contributor to open-source, hub makes it easier to fetch repositories, navigate project pages, fork repos and even submit pull requests, all from the command-line.

### clone your own project
```bash
$ git clone dotfiles
→ git clone git://github.com/YOUR_USER/dotfiles.git
```

### clone another project
```bash
$ git clone github/hub
→ git clone git://github.com/github/hub.git
```

### open the current project's issues page
```bash
$ git browse -- issues
→ open https://github.com/github/hub/issues
```

### open another project's wiki
```bash
$ git browse mojombo/jekyll wiki
→ open https://github.com/mojombo/jekyll/wiki
```

### Example workflow for contributing to a project:
```bash
$ git clone github/hub
$ cd hub
```

### create a topic branch
```bash
$ git checkout -b feature
  ( making changes ... )
$ git commit -m "done with feature"
```

### It's time to fork the repo!
```bash
$ git fork
→ (forking repo on GitHub...)
→ git remote add YOUR_USER git://github.com/YOUR_USER/hub.git
```

### push the changes to your new remote
```bash
$ git push YOUR_USER feature
```

### open a pull request for the topic branch you've just pushed
```bash
$ git pull-request
→ (opens a text editor for your pull request message)
```

## As an open-source maintainer
Maintaining a project is easier when you can easily fetch from other forks, review pull requests and cherry-pick URLs. You can even create a new repo for your next thing.

### fetch from multiple trusted forks, even if they don't yet exist as remotes
```bash
$ git fetch mislav,cehoffman
→ git remote add mislav git://github.com/mislav/hub.git
→ git remote add cehoffman git://github.com/cehoffman/hub.git
→ git fetch --multiple mislav cehoffman
```

### check out a pull request for review
```bash
$ git checkout https://github.com/github/hub/pull/134
→ (creates a new branch with the contents of the pull request)
```

### directly apply all commits from a pull request to the current branch
```bash
$ git am -3 https://github.com/github/hub/pull/134
```

### cherry-pick a GitHub URL
```bash
$ git cherry-pick https://github.com/xoebus/hub/commit/177eeb8
```

### open the GitHub compare view between two releases
```bash
$ git compare v0.9..v1.0
```

### put compare URL for a topic branch to clipboard
```bash
$ git compare -u feature | pbcopy
```

### create a repo for a new project
```bash
$ git init
$ git add . && git commit -m "It begins."
$ git create -d "My new thing"
→ (creates a new project on GitHub with the name of current directory)
$ git push origin master
Using GitHub for work
Save time at work by opening pull requests for code reviews and pushing to multiple remotes at once. Even GitHub Enterprise is supported.
```

### whitelist your GitHub Enterprise hostname
```bash
$ git config --global --add hub.host my.example.org
```

### open a pull request using a message generated from script, then put its URL to the clipboard
```bash
$ git push origin feature
$ git pull-request -c -F prepared-message.md
→ (URL ready for pasting in a chat room)
```

### push to multiple remotes
```bash
$ git push production,staging
```