source $(dirname $0)/common.sh

get_with_token "/api/v1/courses" $1
