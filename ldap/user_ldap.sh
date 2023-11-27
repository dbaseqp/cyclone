USER=$1
PASSWORD=$2

sed -e "s/{USER}/$USER/" \
    -e "s/{PASSWORD}/$PASSWORD/" \
    ./ldap/user_template.ldif > ./ldap/"$USER"_user_add.ldif
ldapadd -w password -H ldap://localhost:389 -D 'cn=admin,dc=sdc,dc=cpp' -f ./ldap/"$USER"_user_add.ldif
rm ./ldap/"$USER"_user_add.ldif