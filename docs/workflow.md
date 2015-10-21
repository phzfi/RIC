# Workflow

## Summary

- Create new branch (not fork)
- Add commits to your branch
- Open a Pull Request
- Someone else reviews the changes and merges them to the master branch

Remember: Everything in the <code>master</code> branch is always deployable!

## Creating a new branch

Changes you make on a branch don't affect the <code>master</code> branch, so you're free to experiment.

To create a new branch:

	git checkout -b branch-name

Commit your changes and push them to the branch:

	git push origin branch-name

Use descriptive names for the branches, lowercase letters and hyphens.

## Pull requests

<<<<<<< HEAD

=======
### Open a pull request

From the dropdown menu select the branch to which you have committed your changes.

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/02-branch-ahead.png)

Select the base repository. Pull request will be opened for this repository. Currently we are using kmort89/IRC which is forked from phzfi/RIC. 

**NOTE: select kmort89/IRC as the base repository**

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/03a-base-phz.png)

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/03b-base-selection.png)

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/03c-base-master.png)

You can review your changes before submitting the pull request.

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/05-open-review.png)

### Merge and close

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/04a-pull-requests.png)

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/04b-pull-requests.png)

Review the contents of the pull request:

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/04c-pull-review.png)

Merge pull request to <code>master</code>

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/04d-pull-merge.png)

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/04e-pull-confirm.png)

After successfully merging you can delete the branch:

![](https://raw.githubusercontent.com/kmort89/RIC/master/docs/images/workflow/04f-pull-delete.png)
>>>>>>> master

## Resources

Highly recommended reading:

https://guides.github.com/introduction/flow/

https://github.com/Kunena/Kunena-Forum/wiki/Create-a-new-branch-with-git-and-manage-branches
