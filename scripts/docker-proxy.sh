docker-proxy () {
  args=$*
  hostname=`ruby -e "match = ARGV[0].match(/(?<=proxy-hostname=)([^\s]+)/); puts match.captures[0] unless match.nil?" "a ${args} a"`

  if [[ -z "${hostname// }" ]]
  then
    echo "No proxy-hostname label"
  else
    hoststring="10.0.0.100  ${hostname}"

    if grep -Fxq "${hoststring}" "/etc/hosts" 
    then
      echo "${hostname} already in /etc/hosts\n"
    else
      echo "adding ${hostname} to /etc/hosts\n"
      echo "${hoststring}" | sudo tee -a /etc/hosts
    fi

    echo "Starting docker container\n"

    eval "docker run ${args}"
  fi
}
