The Git Flow workflow is a branching model designed for efficient collaboration and release management in software development projects. It provides a structured approach to organizing branches and streamlines, the process of feature development, bug fixing, and release management.


## Branches

The Git Flow workflow defines several types of branches:



* Main Branches:
    * main: Represents the main branch and always reflects the production-ready state of the project.
    * develop: Serves as the integration branch where feature branches are merged.
* Supporting Branches:
    * feature/*: Used for developing new features or enhancements. These branches branch off from the develop branch and are merged back into it.
    * hotfix/*: Created to address critical bugs or issues in the production code. These branches branch off from the main branch and are merged back into both main and develop.
    * bugfix/*: Created to address bugs or issues in the development code. These branches branch off from the develop branch and are merged back into develop as if they were feature branches.
    * release/*: Used for preparing and stabilizing releases. These branches branch off from the develop branch and are merged into both main and develop when ready.


## Workflow

The typical Git Flow workflow consists of the following steps:



1. Creating Feature Branches:
    * Identify a new feature to be implemented.
    * Create a new branch using the naming convention feature/&lt;feature-name>.
    * Implement and test the feature on the feature branch.
2. Merging Feature Branches:
    * Once the feature is complete, merge the feature branch back into the develop branch.
    * Review the changes and ensure the integration is successful.
    * The merge will trigger a deployment of the product to the development environment. Make the needed QA tests.
    * Merge the branch as Squash & Merge, remove any automatically generated description, and define the title as denied in the [conventional commit](https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional).
3. Releasing Versions:
    * Create a release branch from the develop branch using the naming convention release/&lt;version>.
    * Perform final testing, bug fixes, and documentation updates on the release branch.
        1. Ideally, the version will be automatically updated, but we may need to update it manually at the first stages of the project.
    * Create a Pull Request to merge into main
        2. This will trigger the deployment of an ephemeral integration environment from ArgoCD via [this](https://github.com/konstellation-io/konstellation-infrastructure/blob/main/argo-cd/applicationsets/int-cluster-integration-applications.yaml) ApplicationSet.
    * Merge the release branch into main with the Merge Commit strategy.
    * The merge to main will create a new pull request to develop, with an increased patch version.
    * Tag the release commit with the corresponding version number.
    * Once merged, merge the latest changes from the main branch, back to the develop branch with the Merge Commit strategy.
4. Hotfixes:
    * If critical issues arise in the production code, create a hotfix branch from the main branch using the naming convention hotfix/&lt;description>.
    * Implement the necessary fixes on the hotfix branch.
    * Merge the hotfix branch back into both main and develop branches.
    * Tag the hotfix commit with an appropriate version number.
    * Once merged, merge the latest changes from the main branch, back to the develop branch with the Merge Commit strategy.
5. Bug Fixes:
    * If bugs or issues are identified in the development code, create a bugfix branch from the develop branch using the naming convention bugfix/&lt;description>.
    * Implement the necessary fixes on the bugfix branch.
    * Merge the bugfix branch back into the develop branch as if it were a feature branch. Merge the branch as Squash & Merge, remove any automatically generated description, and define the title as denied in the [conventional commit](https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional).
