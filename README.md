### Overview
This Go application utilizes the GitHub API to gather insights into the pull requests of a specified repository. It generates a summary report of opened, closed, and in-progress pull requests within the last week, which can be sent for project tracking purposes.

### Example Report
Here's an example of the email summary report format that the application generates:

**From:** sender@example.com
**To:** recipient@example.com
**Subject:** Weekly Pull Request Summary
**Body:**
```
Hello Team,

Here's the summary of pull requests activity in the last week for the repository "owner/repository-name":

- Opened Pull Requests:
  1. #123: "Title of PR 1" by User1
  2. #124: "Title of PR 2" by User2

- Closed Pull Requests:
  1. #120: "Title of PR 3" by User3
  2. #121: "Title of PR 4" by User4

- In-Progress Pull Requests:
  1. #122: "Title of PR 5" by User5

Please review and take necessary actions.

Best regards,
Your Name
```

