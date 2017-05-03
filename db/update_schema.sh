SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "CREATE DATABASE IF NOT EXISTS slackify;" > commands
echo "USE slackify;" > commands
for f in $SCRIPT_DIR/slackify_*.sql; do
    echo "source $f;" >> commands
done

## TODO: read schema version number from table

mysql -u$MYSQL_USER -p < commands

rm commands
