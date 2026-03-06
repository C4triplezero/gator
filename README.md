Requires Postgres and Go installed to run
can be installed from https://github.com/C4triplezero/gator using go install
create a file called ".gatorconfig.json" in your home directory containing this:
{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable","current_user_name":"c4"}
use the register command to create a user | arument: name
use the login command to switch between users | argument: name
use the addfeed command to added a feed to be scrapped from | arguments: name of post and url
use the agg command to collect posts from your feeds | argument: time between collections
use the browse command to view collected posts optional | argument: amount of posts to view (defaults to 2)