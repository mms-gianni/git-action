# git-action
Manage your github actions with git cli


![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mms-gianni/git-action)
![GitHub top language](https://img.shields.io/github/languages/top/mms-gianni/git-action)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/mms-gianni/git-action/Upload%20Release%20Asset)
![GitHub MIT license](https://img.shields.io/github/license/mms-gianni/git-action)
![Swiss made](https://img.shields.io/badge/swiss%20made-100%25-red)


## Why
- Simplify access to your workflows 
- Keep fokus on your work


![Screenshot](docs/img/screenshot-actions.png?raw=true "Screenshot")

## Installation
Generate a token here : https://github.com/settings/tokens (You need to be loged in)

To export the Github username and organisation is optional. 
### Mac
```
echo 'export GITHUB_TOKEN="XXXXXXXXXXXXXXXXXXXXXXX"' >> ~/.zshrc
echo 'export GITHUB_USERNAME="change-me-to-your-username"' >> ~/.zshrc
echo 'export GITHUB_ORGANISATIONS="klustair,kubernetes"' >> ~/.zshrc
curl https://raw.githubusercontent.com/mms-gianni/git-action/master/cmd/git-action/git-action.mac.64bit -o /usr/local/bin/git-action
chmod +x /usr/local/bin/git-action
```

### Linux 
```
echo 'export GITHUB_TOKEN="XXXXXXXXXXXXXXXXXXXXXXX"' >> ~/.bashrc
echo 'export GITHUB_USERNAME="change-me-to-your-username"' >> ~/.bashrc
echo 'export GITHUB_ORGANISATIONS="klustair,kubernetes"' >> ~/.bashrc
curl https://raw.githubusercontent.com/mms-gianni/git-action/master/cmd/git-action/git-action.linux.64bit -o /usr/local/bin/git-action
chmod +x /usr/local/bin/git-action
```

### Windows
Windows is not implemented yet. But I'm working on it. Pullrequests wellcome. 

## Quick start

### Show list of workflow runs
```
cd /path/to/your/repo
git action status
```

### Trigger a workflow
```
cd /path/to/your/repo
git action run
```

See help for all flags and options