# Milestone Editor

## Installation
* `git clone git@github.com:quentin-sommer/github-milestones-editor.git`
* Create a [token](https://github.com/settings/tokens) for the account you wish to automate and put it in accessToken.txt

## Usage

```
//when creating milestones
./main -title=title -desc=desc -date=25-01-2017 -mask="my-mask*"

//when deleting milestones
./main -remove -title=title -mask="my-mask*"

// you can also use it on all your repositories
./main -title=title -desc=desc -date=25-01-2017
```
:bangbang: date format : yyyy-dd-mm
