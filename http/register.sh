set -x

while test $# -gt 0; do
  case "$1" in
    -name)
		shift
		name=$1
		shift
		;;
	-email)
		shift
		email=$1
		shift
		;;
	-password)
		shift
		password=$1
		shift
		;;
	*)
      	break
      	;;
  esac
done

name=${name:-JohnDoe}
email=${email:-johndoe@mail.com}
password=${password:-foobar1234}

curl -v -X POST \
  -H 'Connection: close' \
  -H 'Content-Type: application/json' \
  -d '{"name": "'$name'", "email": "'$email'", "password": "'$password'"}' \
  http://localhost:8080/api/v1/register
